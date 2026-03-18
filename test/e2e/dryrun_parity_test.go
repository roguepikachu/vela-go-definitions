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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// skipParityTests lists test files that cannot be used with dry-run
// (e.g., they reference external resources or require multi-cluster setup).
var skipParityTests = map[string]string{
	"deploy-cloud-resource.yaml":    "requires alibaba-rds component and multi-cluster setup",
	"share-cloud-resource.yaml":     "requires alibaba-rds component and multi-cluster setup",
	"generate-jdbc-connection.yaml": "requires alibaba-rds component",
	"apply-terraform-config.yaml":   "requires Alibaba Cloud credentials and terraform provider",
	"apply-terraform-provider.yaml": "requires Alibaba Cloud credentials",
}

var _ = Describe("Dry-Run Parity Tests", Label("parity"), func() {
	baselineDir := os.Getenv("CUE_BASELINE_DIR")

	Context("when comparing defkit dry-run output against CUE baseline", func() {
		BeforeEach(func() {
			if baselineDir == "" {
				Skip("CUE_BASELINE_DIR not set — skipping parity tests. " +
					"Generate baselines with: make generate-baseline BASELINE_DIR=/tmp/cue-baseline")
			}
		})

		// Test each definition type subdirectory
		paritySubdirs := []struct {
			name   string
			subdir string
		}{
			{name: "components", subdir: "applications/components"},
			{name: "traits", subdir: "applications/trait"},
			{name: "policies", subdir: "applications/policies"},
			{name: "workflowsteps", subdir: "applications/workflowsteps"},
		}

		for _, ps := range paritySubdirs {
			ps := ps
			When(fmt.Sprintf("testing %s parity", ps.name), func() {
				testDataPath := filepath.Join(getTestDataPath(), ps.subdir)

				for _, file := range func() []string {
					f, _ := listYAMLFiles(testDataPath)
					return f
				}() {
					file := file

					It(fmt.Sprintf("should match CUE baseline for %s", filepath.Base(file)), func() {
						baseName := filepath.Base(file)

						if reason, ok := skipParityTests[baseName]; ok {
							Skip(fmt.Sprintf("Skipping: %s", reason))
						}

						// Look for baseline file
						baselineFile := filepath.Join(baselineDir, strings.TrimSuffix(baseName, filepath.Ext(baseName))+".yaml")
						if _, err := os.Stat(baselineFile); os.IsNotExist(err) {
							Skip(fmt.Sprintf("No baseline file found at %s — skipping", baselineFile))
						}

						// Read baseline
						baselineData, err := os.ReadFile(baselineFile)
						Expect(err).NotTo(HaveOccurred(), "Failed to read baseline file %s", baselineFile)

						// Run dry-run with defkit definitions
						actualOutput, err := runVelaDryRun(file)
						Expect(err).NotTo(HaveOccurred(), "vela dry-run failed for %s", file)

						// Parse both into normalized resource maps
						baselineResources, err := parseDryRunResources(string(baselineData))
						Expect(err).NotTo(HaveOccurred(), "Failed to parse baseline resources")

						actualResources, err := parseDryRunResources(actualOutput)
						Expect(err).NotTo(HaveOccurred(), "Failed to parse defkit dry-run resources")

						// Compare
						diffs := compareResources(baselineResources, actualResources)
						if len(diffs) > 0 {
							GinkgoWriter.Printf("\n=== Parity differences for %s ===\n", baseName)
							for _, d := range diffs {
								GinkgoWriter.Printf("  - %s\n", d)
							}
						}
						Expect(diffs).To(BeEmpty(),
							fmt.Sprintf("Defkit output does not match CUE baseline for %s:\n%s",
								baseName, strings.Join(diffs, "\n")))
					})
				}
			})
		}
	})
})
