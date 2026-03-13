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

var _ = Describe("Export2Config WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Export2Config()
			Expect(step.GetName()).To(Equal("export2config"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Export2Config()
			Expect(step.GetDescription()).To(Equal("Export data to specified Kubernetes ConfigMap in your workflow."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Export2Config()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required configName", func() {
				Expect(cueOutput).To(ContainSubstring("configName: string"))
			})

			It("should have optional namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			})

			It("should have required data as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("data: {}"))
			})

			It("should have cluster with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			})
		})

		Describe("Template: kube.#Apply", func() {
			It("should use kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should create a v1 ConfigMap", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "ConfigMap"`))
			})

			It("should set metadata name from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("name: parameter.configName"))
			})

			It("should reference parameter data", func() {
				Expect(cueOutput).To(ContainSubstring("data: parameter.data"))
			})

			It("should pass cluster parameter", func() {
				Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
			})
		})

		Describe("Template: namespace guards", func() {
			It("should use mutually exclusive namespace guards", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] != _|_`))
				Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] == _|_`))
			})

			It("should set namespace to parameter.namespace when set", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})

			It("should set namespace to context.namespace when not set", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should NOT have unconditional namespace assignment", func() {
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
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
