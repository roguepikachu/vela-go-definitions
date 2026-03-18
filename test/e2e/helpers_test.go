/*
Copyright 2025 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/yaml"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

const (
	// Timeout for application to become running
	AppRunningTimeout = 5 * time.Minute
	// Polling interval for status checks
	PollInterval = 5 * time.Second
)

var (
	k8sClient client.Client
)

// initK8sClient initializes the Kubernetes controller-runtime client once.
func initK8sClient() error {
	if k8sClient != nil {
		return nil
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	// Register KubeVela schemes
	_ = v1beta1.AddToScheme(scheme.Scheme)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}

	return nil
}

// readAppFromFile reads an Application from a YAML file (supports multi-doc YAML).
func readAppFromFile(filename string) (*v1beta1.Application, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	docs := strings.Split(string(bs), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		app := &v1beta1.Application{}
		if err = yaml.Unmarshal([]byte(doc), app); err != nil {
			continue
		}

		if app.Kind == "Application" {
			return app, nil
		}
	}

	return nil, fmt.Errorf("no Application found in file %s", filename)
}

// updateAppNamespaceReferences updates namespace references inside Application components
// This is needed for ref-objects type components that reference resources in specific namespaces
func updateAppNamespaceReferences(app *v1beta1.Application, newNamespace string) {
	for i := range app.Spec.Components {
		comp := &app.Spec.Components[i]
		if comp.Type == "ref-objects" && comp.Properties != nil {
			// Parse properties as map
			var props map[string]interface{}
			if err := json.Unmarshal(comp.Properties.Raw, &props); err != nil {
				continue
			}

			// Update namespace in objects array
			if objects, ok := props["objects"].([]interface{}); ok {
				for _, obj := range objects {
					if objMap, ok := obj.(map[string]interface{}); ok {
						// Update namespace field to the new namespace
						if _, hasNs := objMap["namespace"]; hasNs {
							objMap["namespace"] = newNamespace
						}
					}
				}
			}

			// Marshal back to properties
			if newProps, err := json.Marshal(props); err == nil {
				comp.Properties = &runtime.RawExtension{Raw: newProps}
			}
		}
	}
}

// getProjectRoot finds the project root by looking for go.mod.
func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// getTestDataPath returns the path to the test data directory.
func getTestDataPath() string {
	if path := os.Getenv("TESTDATA_PATH"); path != "" {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(getProjectRoot(), path)
	}
	return filepath.Join(getProjectRoot(), "test", "builtin-definition-example")
}

// listYAMLFiles lists all YAML files in a directory.
func listYAMLFiles(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	return files, nil
}

// sanitizeForNamespace creates a DNS-1123 compliant name
func sanitizeForNamespace(name string) string {
	n := strings.ToLower(name)
	n = strings.ReplaceAll(n, "_", "-")
	n = strings.ReplaceAll(n, ".", "-")
	// Keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range n {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	out := strings.Trim(result.String(), "-")
	if len(out) > 30 {
		out = out[:30]
	}
	return strings.Trim(out, "-")
}

// applyPrerequisiteResources applies non-Application resources from a multi-doc YAML file
// This is needed for files like ref-objects.yaml that reference existing Deployment/Service
func applyPrerequisiteResources(ctx context.Context, filePath, namespace string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	docs := strings.Split(string(content), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		// Parse as unstructured to check the kind
		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(doc), obj); err != nil {
			continue
		}

		// Skip Application resources - we'll apply those separately
		if obj.GetKind() == "Application" {
			continue
		}

		// Skip empty objects
		if obj.GetKind() == "" {
			continue
		}

		// Update namespace
		obj.SetNamespace(namespace)

		GinkgoWriter.Printf("Applying prerequisite %s/%s in namespace %s...\n", obj.GetKind(), obj.GetName(), namespace)

		err = k8sClient.Create(ctx, obj)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// Update if already exists
				existing := &unstructured.Unstructured{}
				existing.SetGroupVersionKind(obj.GroupVersionKind())
				if getErr := k8sClient.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: namespace}, existing); getErr == nil {
					obj.SetResourceVersion(existing.GetResourceVersion())
					err = k8sClient.Update(ctx, obj)
				}
			}
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("failed to apply %s/%s: %w", obj.GetKind(), obj.GetName(), err)
			}
		}
	}

	return nil
}

// hasPrerequisiteResources checks if a YAML file contains non-Application resources
func hasPrerequisiteResources(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	docs := strings.Split(string(content), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(doc), obj); err != nil {
			continue
		}

		// If there's any non-Application resource, return true
		if obj.GetKind() != "" && obj.GetKind() != "Application" {
			return true
		}
	}

	return false
}

// getAppFailureDiagnostics gathers detailed diagnostic information when an application fails.
// This includes application status, workflow step details, vela status output, and kubectl describe output.
func getAppFailureDiagnostics(ctx context.Context, appName, namespace string) string {
	currentApp := &v1beta1.Application{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: appName}, currentApp); err != nil {
		return fmt.Sprintf("Failed to get app for diagnostics: %v", err)
	}

	var diagInfo strings.Builder
	diagInfo.WriteString(fmt.Sprintf("\n=== Application %s/%s Failed ===\n", namespace, appName))
	diagInfo.WriteString(fmt.Sprintf("Phase: %s\n", currentApp.Status.Phase))

	// Workflow status details
	if currentApp.Status.Workflow != nil {
		diagInfo.WriteString(fmt.Sprintf("Workflow Mode: %s\n", currentApp.Status.Workflow.Mode))
		diagInfo.WriteString(fmt.Sprintf("Workflow Finished: %v\n", currentApp.Status.Workflow.Finished))
		diagInfo.WriteString(fmt.Sprintf("Workflow Terminated: %v\n", currentApp.Status.Workflow.Terminated))
		diagInfo.WriteString(fmt.Sprintf("Workflow Suspended: %v\n", currentApp.Status.Workflow.Suspend))
		diagInfo.WriteString(fmt.Sprintf("Workflow Message: %s\n", currentApp.Status.Workflow.Message))

		diagInfo.WriteString("\n--- Workflow Steps ---\n")
		for _, step := range currentApp.Status.Workflow.Steps {
			diagInfo.WriteString(fmt.Sprintf("  Step: %s (type: %s)\n", step.Name, step.Type))
			diagInfo.WriteString(fmt.Sprintf("    Phase: %s\n", step.Phase))
			diagInfo.WriteString(fmt.Sprintf("    Message: %s\n", step.Message))
			diagInfo.WriteString(fmt.Sprintf("    Reason: %s\n", step.Reason))
		}
	}

	// Run vela status command
	cmd := exec.Command("vela", "status", appName, "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		diagInfo.WriteString(fmt.Sprintf("\nvela status error: %v\n", err))
	}
	diagInfo.WriteString(fmt.Sprintf("\n--- vela status output ---\n%s\n", string(output)))

	// Describe application via kubectl
	descCmd := exec.Command("kubectl", "describe", "app", appName, "-n", namespace)
	descOutput, err := descCmd.CombinedOutput()
	if err != nil {
		diagInfo.WriteString(fmt.Sprintf("\nkubectl describe error: %v\n", err))
	}
	diagInfo.WriteString(fmt.Sprintf("\n--- kubectl describe app ---\n%s\n", string(descOutput)))

	// List all pods in the namespace
	podsCmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "wide")
	podsOutput, err := podsCmd.CombinedOutput()
	if err != nil {
		diagInfo.WriteString(fmt.Sprintf("\nkubectl get pods error: %v\n", err))
	}
	diagInfo.WriteString(fmt.Sprintf("\n--- Pods in namespace %s ---\n%s\n", namespace, string(podsOutput)))

	return diagInfo.String()
}

// getKanikoPodLogs retrieves logs and details from kaniko pods in a namespace.
// This is specifically used for build-push-image workflow step diagnostics.
func getKanikoPodLogs(namespace string) string {
	var diagInfo strings.Builder

	kanikoPodCmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "name", "--no-headers")
	kanikoPodOutput, err := kanikoPodCmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error listing pods: %v\n", err)
	}

	podNames := strings.Split(strings.TrimSpace(string(kanikoPodOutput)), "\n")
	for _, podName := range podNames {
		if podName == "" {
			continue
		}
		// Only get logs from kaniko pods (created by build-push-image workflow step)
		if !strings.Contains(podName, "kaniko") {
			continue
		}
		// Get logs for kaniko pod
		diagInfo.WriteString(fmt.Sprintf("\n--- Logs from %s ---\n", podName))
		logsCmd := exec.Command("kubectl", "logs", podName, "-n", namespace, "--tail=100")
		logsOutput, err := logsCmd.CombinedOutput()
		if err != nil {
			diagInfo.WriteString(fmt.Sprintf("Error getting logs: %v\n", err))
		} else {
			diagInfo.WriteString(string(logsOutput))
		}

		// Also describe the pod for events and status
		diagInfo.WriteString(fmt.Sprintf("\n--- Describe %s ---\n", podName))
		describePodCmd := exec.Command("kubectl", "describe", podName, "-n", namespace)
		describePodOutput, err := describePodCmd.CombinedOutput()
		if err != nil {
			diagInfo.WriteString(fmt.Sprintf("Error describing pod: %v\n", err))
		} else {
			diagInfo.WriteString(string(describePodOutput))
		}
	}

	if diagInfo.Len() == 0 {
		return "No kaniko pods found in namespace\n"
	}

	return diagInfo.String()
}

// --------------------------------------------------------------------------
// Shared test runner (Phase 1)
// --------------------------------------------------------------------------

// skipWorkflowStepTests lists test files that require external infrastructure
// (cloud providers, terraform, Prometheus, container registries, webhook endpoints)
// and cannot run in a standard CI environment.
var skipWorkflowStepTests = map[string]string{
	"deploy-cloud-resource.yaml":    "requires alibaba-rds component and multi-cluster setup",
	"share-cloud-resource.yaml":     "requires alibaba-rds component and multi-cluster setup",
	"generate-jdbc-connection.yaml": "requires alibaba-rds component",
	"apply-terraform-config.yaml":   "requires Alibaba Cloud credentials and terraform provider",
	"apply-terraform-provider.yaml": "requires Alibaba Cloud credentials",
}

// runDefinitionTest executes a single definition e2e test case.
// It creates an isolated namespace, applies the application, waits for running status,
// optionally validates resource expectations, and cleans up.
func runDefinitionTest(ctx context.Context, file string, skipTests map[string]string) {
	if reason, ok := skipTests[filepath.Base(file)]; ok {
		Skip(fmt.Sprintf("Skipping: %s", reason))
	}

	app, err := readAppFromFile(file)
	Expect(err).NotTo(HaveOccurred(), "Failed to read application from %s", file)

	// Each app has unique name, so namespace based on app name is unique per test
	appNameSanitized := sanitizeForNamespace(app.Name)
	uniqueNs := fmt.Sprintf("e2e-%s", appNameSanitized)

	app.SetNamespace(uniqueNs)

	// Update namespace references inside component properties (e.g., ref-objects)
	updateAppNamespaceReferences(app, uniqueNs)

	// Track test success for cleanup diagnostics
	testPassed := false

	// DeferCleanup runs even on suite timeout - print diagnostics if test didn't pass
	DeferCleanup(func() {
		if !testPassed {
			GinkgoWriter.Printf("\nTest did not complete successfully, gathering diagnostics...\n")
			GinkgoWriter.Printf("%s\n", getAppFailureDiagnostics(ctx, app.Name, uniqueNs))
			// Print kaniko pod logs for build-push-image workflow step
			if filepath.Base(file) == "build-push-image.yaml" {
				GinkgoWriter.Printf("\n--- Kaniko Pod Diagnostics ---\n")
				GinkgoWriter.Printf("%s\n", getKanikoPodLogs(uniqueNs))
			}
		}
		// Clean up namespace after test
		GinkgoWriter.Printf("Deleting namespace %s...\n", uniqueNs)
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: uniqueNs}}
		_ = k8sClient.Delete(ctx, ns)
	})

	GinkgoWriter.Printf("Creating namespace %s...\n", uniqueNs)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: uniqueNs}}
	err = k8sClient.Create(ctx, ns)
	if err != nil && !errors.IsAlreadyExists(err) {
		Expect(err).NotTo(HaveOccurred(), "Failed to create namespace")
	}

	// Ensure clean slate - delete app if exists
	GinkgoWriter.Printf("Cleaning up any existing application %s/%s...\n", uniqueNs, app.Name)
	_ = k8sClient.Delete(ctx, app)

	// Wait for deletion
	Eventually(func() bool {
		err := k8sClient.Get(ctx, types.NamespacedName{Namespace: uniqueNs, Name: app.Name}, &v1beta1.Application{})
		return errors.IsNotFound(err)
	}, 30*time.Second, 2*time.Second).Should(BeTrue(),
		fmt.Sprintf("Application %s should be fully deleted before test", app.Name))

	// Check if this file has prerequisite resources
	if hasPrerequisiteResources(file) {
		GinkgoWriter.Printf("Applying prerequisite resources from %s...\n", filepath.Base(file))
		err = applyPrerequisiteResources(ctx, file, uniqueNs)
		Expect(err).NotTo(HaveOccurred(), "Failed to apply prerequisite resources")
		time.Sleep(2 * time.Second)
	}

	// Apply application
	GinkgoWriter.Printf("Applying application %s/%s...\n", uniqueNs, app.Name)
	Expect(k8sClient.Create(ctx, app)).Should(Succeed())

	// Wait for running status
	Eventually(func(g Gomega) {
		currentApp := &v1beta1.Application{}
		g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: uniqueNs, Name: app.Name}, currentApp)).Should(Succeed())
		GinkgoWriter.Printf("Application %s status: %s\n", app.Name, currentApp.Status.Phase)
		g.Expect(string(currentApp.Status.Phase)).Should(Equal("running"))
	}, AppRunningTimeout, PollInterval).Should(Succeed())

	// Validate resource expectations if companion .expect.yaml exists
	expectations := loadExpectations(file)
	if len(expectations) > 0 {
		GinkgoWriter.Printf("Validating %d resource expectation(s)...\n", len(expectations))
		validateResourceExpectations(ctx, uniqueNs, expectations)
	}

	testPassed = true
	GinkgoWriter.Printf("PASS %s\n", filepath.Base(file))
}

// --------------------------------------------------------------------------
// Resource expectation validation (Phase 2)
// --------------------------------------------------------------------------

// ResourceExpectation describes expected state of a K8s resource after an Application is running.
type ResourceExpectation struct {
	APIVersion string                 `yaml:"apiVersion" json:"apiVersion"`
	Kind       string                 `yaml:"kind" json:"kind"`
	Name       string                 `yaml:"name" json:"name"`
	Fields     map[string]interface{} `yaml:"fields" json:"fields"`
}

// ExpectationFile is the top-level structure of a .expect.yaml file.
type ExpectationFile struct {
	Expectations []ResourceExpectation `yaml:"expectations" json:"expectations"`
}

// loadExpectations looks for a .expect.yaml file in the expectations/ directory
// that mirrors the applications/ directory structure.
// For example, given .../builtin-definition-example/applications/components/webservice.yaml,
// it looks for .../builtin-definition-example/expectations/components/webservice.expect.yaml.
// Returns nil if no expectation file exists.
func loadExpectations(appYAMLPath string) []ResourceExpectation {
	// appYAMLPath: .../builtin-definition-example/applications/<type>/<name>.yaml
	// expectPath:  .../builtin-definition-example/expectations/<type>/<name>.expect.yaml
	dir := filepath.Dir(appYAMLPath)             // .../applications/components
	subdir := filepath.Base(dir)                 // components
	testDataRoot := filepath.Dir(filepath.Dir(dir)) // .../builtin-definition-example
	baseName := filepath.Base(appYAMLPath)       // webservice.yaml
	ext := filepath.Ext(baseName)                // .yaml
	nameNoExt := strings.TrimSuffix(baseName, ext) // webservice

	expectPath := filepath.Join(testDataRoot, "expectations", subdir, nameNoExt+".expect.yaml")

	data, err := os.ReadFile(expectPath)
	if err != nil {
		return nil // No expectation file — that's fine
	}

	var ef ExpectationFile
	if err := yaml.Unmarshal(data, &ef); err != nil {
		GinkgoWriter.Printf("Warning: failed to parse %s: %v\n", expectPath, err)
		return nil
	}

	return ef.Expectations
}

// parseGVK parses an apiVersion and kind into a GroupVersionKind.
func parseGVK(apiVersion, kind string) schema.GroupVersionKind {
	parts := strings.SplitN(apiVersion, "/", 2)
	if len(parts) == 1 {
		// core group, e.g. "v1"
		return schema.GroupVersionKind{Group: "", Version: parts[0], Kind: kind}
	}
	return schema.GroupVersionKind{Group: parts[0], Version: parts[1], Kind: kind}
}

// validateResourceExpectations fetches each expected resource and validates its fields.
func validateResourceExpectations(ctx context.Context, namespace string, expectations []ResourceExpectation) {
	for _, exp := range expectations {
		GinkgoWriter.Printf("  Checking %s/%s %s...\n", exp.APIVersion, exp.Kind, exp.Name)

		obj := &unstructured.Unstructured{}
		obj.SetGroupVersionKind(parseGVK(exp.APIVersion, exp.Kind))

		// Fetch the resource — retry briefly in case of propagation delay
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: exp.Name}, obj)
		}, 30*time.Second, 2*time.Second).Should(Succeed(),
			fmt.Sprintf("Expected %s/%s %q to exist in namespace %s", exp.APIVersion, exp.Kind, exp.Name, namespace))

		// Validate each field path
		for path, expectedValue := range exp.Fields {
			actual, err := getNestedValue(obj.Object, path)
			Expect(err).NotTo(HaveOccurred(), "Failed to resolve path %q in %s/%s %s", path, exp.APIVersion, exp.Kind, exp.Name)

			// Normalize numbers for comparison (YAML/JSON may parse as float64 or int64)
			assertValuesEqual(path, expectedValue, actual,
				fmt.Sprintf("%s/%s %s", exp.APIVersion, exp.Kind, exp.Name))
		}
	}
}

// arrayIndexPattern matches path segments like "containers[0]"
var arrayIndexPattern = regexp.MustCompile(`^(.+)\[(\d+)\]$`)

// getNestedValue walks a dot-path with optional array indexing into an unstructured object.
// Examples: "spec.replicas", "spec.template.spec.containers[0].image"
func getNestedValue(obj map[string]interface{}, path string) (interface{}, error) {
	segments := splitDotPath(path)
	var current interface{} = obj

	for _, seg := range segments {
		if current == nil {
			return nil, fmt.Errorf("nil value at segment %q in path %q", seg, path)
		}

		// Check for array index: "containers[0]"
		if m := arrayIndexPattern.FindStringSubmatch(seg); m != nil {
			fieldName := m[1]
			idx, _ := strconv.Atoi(m[2])

			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected map at %q, got %T", seg, current)
			}
			arr, ok := currentMap[fieldName]
			if !ok {
				return nil, fmt.Errorf("field %q not found", fieldName)
			}
			slice, ok := arr.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected array at %q, got %T", fieldName, arr)
			}
			if idx >= len(slice) {
				return nil, fmt.Errorf("index %d out of bounds (len=%d) at %q", idx, len(slice), fieldName)
			}
			current = slice[idx]
		} else {
			// Simple field access
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected map at %q, got %T", seg, current)
			}
			val, ok := currentMap[seg]
			if !ok {
				return nil, fmt.Errorf("field %q not found", seg)
			}
			current = val
		}
	}

	return current, nil
}

// splitDotPath splits a dot-path while respecting array indices.
// "spec.template.spec.containers[0].image" -> ["spec", "template", "spec", "containers[0]", "image"]
func splitDotPath(path string) []string {
	return strings.Split(path, ".")
}

// assertValuesEqual compares expected and actual values with type normalization.
func assertValuesEqual(path string, expected, actual interface{}, resourceDesc string) {
	// Normalize both sides for comparison
	expected = normalizeValue(expected)
	actual = normalizeValue(actual)

	Expect(reflect.DeepEqual(expected, actual)).To(BeTrue(),
		fmt.Sprintf("Field %q in %s:\n  expected: %v (%T)\n  actual:   %v (%T)",
			path, resourceDesc, expected, expected, actual, actual))
}

// normalizeValue normalizes a value for comparison.
// JSON/YAML can represent numbers as float64, int64, or int — this normalizes them.
func normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case float64:
		// If it's a whole number, convert to int64 for comparison
		if val == float64(int64(val)) {
			return int64(val)
		}
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = normalizeValue(item)
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{}, len(val))
		for k, item := range val {
			result[k] = normalizeValue(item)
		}
		return result
	default:
		return v
	}
}

// --------------------------------------------------------------------------
// Dry-run parity helpers (Phase 3)
// --------------------------------------------------------------------------

// runVelaDryRun executes `vela dry-run -f <appFile>` and returns the rendered output.
func runVelaDryRun(appFile string) (string, error) {
	cmd := exec.Command("vela", "dry-run", "-f", appFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("vela dry-run failed for %s: %w\nOutput: %s", appFile, err, string(output))
	}
	return string(output), nil
}

// parseDryRunResources splits vela dry-run output into individual YAML documents
// and returns them as normalized unstructured objects keyed by "Kind/Name".
func parseDryRunResources(dryRunOutput string) (map[string]map[string]interface{}, error) {
	resources := make(map[string]map[string]interface{})

	docs := strings.Split(dryRunOutput, "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		var obj map[string]interface{}
		if err := yaml.Unmarshal([]byte(doc), &obj); err != nil {
			continue
		}

		kind, _ := obj["kind"].(string)
		metadata, _ := obj["metadata"].(map[string]interface{})
		name := ""
		if metadata != nil {
			name, _ = metadata["name"].(string)
		}

		if kind == "" {
			continue
		}

		key := fmt.Sprintf("%s/%s", kind, name)
		// Strip volatile metadata fields for comparison
		normalizeForParity(obj)
		resources[key] = obj
	}

	return resources, nil
}

// normalizeForParity removes fields that are expected to differ between CUE and defkit
// but don't represent functional differences.
func normalizeForParity(obj map[string]interface{}) {
	// Remove volatile metadata fields
	if metadata, ok := obj["metadata"].(map[string]interface{}); ok {
		delete(metadata, "resourceVersion")
		delete(metadata, "uid")
		delete(metadata, "creationTimestamp")
		delete(metadata, "generation")
		delete(metadata, "managedFields")
		delete(metadata, "selfLink")

		// Remove empty annotations/labels
		if ann, ok := metadata["annotations"].(map[string]interface{}); ok && len(ann) == 0 {
			delete(metadata, "annotations")
		}
		if labels, ok := metadata["labels"].(map[string]interface{}); ok && len(labels) == 0 {
			delete(metadata, "labels")
		}
	}

	// Remove status (not part of rendered output comparison)
	delete(obj, "status")

	// Recursively normalize nested objects
	for _, v := range obj {
		switch val := v.(type) {
		case map[string]interface{}:
			normalizeForParity(val)
		case []interface{}:
			for _, item := range val {
				if m, ok := item.(map[string]interface{}); ok {
					normalizeForParity(m)
				}
			}
		}
	}
}

// compareResources compares two sets of parsed dry-run resources.
// Returns a list of difference descriptions, or nil if they match.
func compareResources(baseline, actual map[string]map[string]interface{}) []string {
	var diffs []string

	// Check for resources in baseline but missing from actual
	for key := range baseline {
		if _, ok := actual[key]; !ok {
			diffs = append(diffs, fmt.Sprintf("Resource %q present in CUE baseline but missing from defkit output", key))
		}
	}

	// Check for resources in actual but missing from baseline
	for key := range actual {
		if _, ok := baseline[key]; !ok {
			diffs = append(diffs, fmt.Sprintf("Resource %q present in defkit output but missing from CUE baseline", key))
		}
	}

	// Compare matching resources
	for key, baseObj := range baseline {
		actualObj, ok := actual[key]
		if !ok {
			continue
		}

		baseJSON, _ := json.Marshal(normalizeValue(baseObj))
		actualJSON, _ := json.Marshal(normalizeValue(actualObj))

		if string(baseJSON) != string(actualJSON) {
			diffs = append(diffs, fmt.Sprintf("Resource %q differs between CUE baseline and defkit output", key))
		}
	}

	return diffs
}
