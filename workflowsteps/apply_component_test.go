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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("ApplyComponent WorkflowStep", func() {
	It("should have correct name and description", func() {
		step := workflowsteps.ApplyComponent()
		Expect(step.GetName()).To(Equal("apply-component"))
		Expect(step.GetDescription()).To(Equal("Apply a specific component and its corresponding traits in application"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyComponent()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, scope, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
			Expect(cueOutput).To(ContainSubstring(`"apply-component": {`))
		})

		It("should declare component, cluster, and namespace parameters", func() {
			Expect(cueOutput).To(ContainSubstring("component: string"))
			Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			Expect(cueOutput).To(ContainSubstring(`namespace: *"" | string`))
		})

		It("should have no imports or template actions", func() {
			Expect(cueOutput).NotTo(ContainSubstring("import"))
			Expect(cueOutput).NotTo(ContainSubstring("kube."))
			Expect(cueOutput).NotTo(ContainSubstring("builtin."))
		})
	})
})
