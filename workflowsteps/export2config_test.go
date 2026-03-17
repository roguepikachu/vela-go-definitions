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
	It("should have correct metadata", func() {
		step := workflowsteps.Export2Config()
		Expect(step.GetName()).To(Equal("export2config"))
		Expect(step.GetDescription()).To(Equal("Export data to specified Kubernetes ConfigMap in your workflow."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Export2Config()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
		})

		It("should import vela/kube", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
		})

		It("should declare all parameters with correct types", func() {
			Expect(cueOutput).To(ContainSubstring("configName: string"))
			Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			Expect(cueOutput).To(ContainSubstring("data: {}"))
			Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
		})

		It("should generate kube.#Apply with ConfigMap resource", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "ConfigMap"`))
			Expect(cueOutput).To(ContainSubstring("name: parameter.configName"))
			Expect(cueOutput).To(ContainSubstring("data: parameter.data"))
			Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
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

		It("should have exactly one kube.#Apply", func() {
			count := strings.Count(cueOutput, "kube.#Apply & {")
			Expect(count).To(Equal(1))
		})
	})
})
