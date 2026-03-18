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

var _ = Describe("Deploy WorkflowStep", func() {
	It("should have correct name and description", func() {
		step := workflowsteps.Deploy()
		Expect(step.GetName()).To(Equal("deploy"))
		Expect(step.GetDescription()).To(Equal("A powerful and unified deploy step for components multi-cluster delivery with policies."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Deploy()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and scope", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
		})

		It("should import vela/multicluster and vela/builtin", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/multicluster"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should declare all parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("auto: *true | bool"))
			Expect(cueOutput).To(ContainSubstring("policies: *[] | [...string]"))
			Expect(cueOutput).To(ContainSubstring("parallelism: *5 | int"))
			Expect(cueOutput).To(ContainSubstring("ignoreTerraformComponent: *true | bool"))
		})

		It("should conditionally suspend when auto is false with correct message", func() {
			Expect(cueOutput).To(ContainSubstring("if parameter.auto == false"))
			Expect(cueOutput).To(ContainSubstring("builtin.#Suspend & {"))
			Expect(cueOutput).To(ContainSubstring(`"Waiting approval to the deploy step \"\(context.stepName)\""`))
		})

		It("should invoke multicluster.#Deploy with all required parameters", func() {
			Expect(cueOutput).To(ContainSubstring("multicluster.#Deploy & {"))
			Expect(cueOutput).To(ContainSubstring("policies: parameter.policies"))
			Expect(cueOutput).To(ContainSubstring("parallelism: parameter.parallelism"))
			Expect(cueOutput).To(ContainSubstring("ignoreTerraformComponent: parameter.ignoreTerraformComponent"))
		})

		It("should have exactly one multicluster.#Deploy and one builtin.#Suspend", func() {
			Expect(strings.Count(cueOutput, "multicluster.#Deploy & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#Suspend & {")).To(Equal(1))
		})
	})
})
