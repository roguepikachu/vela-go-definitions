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

package workflowsteps_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("Export2Secret WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.Export2Secret()
		Expect(step.GetName()).To(Equal("export2secret"))
		Expect(step.GetDescription()).To(Equal("Export data to Kubernetes Secret in your workflow."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Export2Secret()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
		})

		It("should import vela/kube, encoding/base64, and encoding/json", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
		})

		It("should declare all parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("secretName: string"))
			Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			Expect(cueOutput).To(ContainSubstring("type?: string"))
			Expect(cueOutput).To(ContainSubstring("data: {}"))
			Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			Expect(cueOutput).To(ContainSubstring(`kind: *"generic" | "docker-registry"`))
			Expect(cueOutput).To(ContainSubstring("dockerRegistry?: {"))
			Expect(cueOutput).To(ContainSubstring("username: string"))
			Expect(cueOutput).To(ContainSubstring("password: string"))
			Expect(cueOutput).To(ContainSubstring(`server: *"https://index.docker.io/v1/" | string`))
		})

		It("should wrap template in secret block with data helper", func() {
			Expect(cueOutput).To(ContainSubstring("secret: {"))
			Expect(cueOutput).To(ContainSubstring("data: *parameter.data | {}"))
		})

		It("should use mutually exclusive namespace guards", func() {
			Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] != _|_`))
			Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] == _|_`))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))

			lines := strings.Split(cueOutput, "\n")
			for i, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "namespace: context.namespace" {
					found := false
					for j := i - 1; j >= 0; j-- {
						prev := strings.TrimSpace(lines[j])
						if prev == "" {
							continue
						}
						if strings.Contains(prev, "if ") {
							found = true
						}
						break
					}
					Expect(found).To(BeTrue(), "namespace: context.namespace should be inside an if block")
				}
			}
		})

		It("should create a v1 Secret via kube.#Apply with correct fields", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Secret"`))
			Expect(cueOutput).To(ContainSubstring("name: parameter.secretName"))
			Expect(cueOutput).To(ContainSubstring("stringData: data"))
			Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
		})

		It("should conditionally set type based on kind and type parameter", func() {
			Expect(cueOutput).To(ContainSubstring(`parameter["type"] == _|_ && parameter.kind == "docker-registry"`))
			Expect(cueOutput).To(ContainSubstring(`type: "kubernetes.io/dockerconfigjson"`))
			Expect(cueOutput).To(ContainSubstring(`parameter["type"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("type: parameter.type"))
		})

		It("should handle docker registry mode with base64-encoded auth", func() {
			Expect(cueOutput).To(ContainSubstring(`parameter.kind == "docker-registry" && parameter["dockerRegistry"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("registryData:"))
			Expect(cueOutput).To(ContainSubstring("auths:"))
			Expect(cueOutput).To(ContainSubstring(`"\(parameter.dockerRegistry.server)"`))
			Expect(cueOutput).To(ContainSubstring("username: parameter.dockerRegistry.username"))
			Expect(cueOutput).To(ContainSubstring("password: parameter.dockerRegistry.password"))
			Expect(cueOutput).To(ContainSubstring("base64.Encode(null,"))
			Expect(cueOutput).To(ContainSubstring(`".dockerconfigjson": json.Marshal(registryData)`))
		})

		It("should have exactly one kube.#Apply and one secret block", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "\tsecret: {")).To(Equal(1))
		})
	})
})
