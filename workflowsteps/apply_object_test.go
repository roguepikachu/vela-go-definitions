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

var _ = Describe("ApplyObject WorkflowStep", func() {
	It("should have correct name and description", func() {
		step := workflowsteps.ApplyObject()
		Expect(step.GetName()).To(Equal("apply-object"))
		Expect(step.GetDescription()).To(Equal("Apply raw kubernetes objects for your workflow steps"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyObject()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			Expect(cueOutput).To(ContainSubstring(`"apply-object": {`))
		})

		It("should import vela/kube", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
		})

		It("should declare value and cluster parameters", func() {
			Expect(cueOutput).To(ContainSubstring("value: {...}"))
			Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
		})

		It("should generate template with exactly one kube.#Apply passing full parameter object", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			count := strings.Count(cueOutput, "kube.#Apply & {")
			Expect(count).To(Equal(1))
		})
	})
})
