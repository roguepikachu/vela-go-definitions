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

var _ = Describe("ApplyTerraformConfig WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.ApplyTerraformConfig()
		Expect(step.GetName()).To(Equal("apply-terraform-config"))
		Expect(step.GetDescription()).To(Equal("Apply terraform configuration in the step"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyTerraformConfig()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, alias, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Terraform"`))
			Expect(cueOutput).To(ContainSubstring(`alias: ""`))
			Expect(cueOutput).To(ContainSubstring(`"apply-terraform-config": {`))
		})

		It("should import vela/kube and vela/builtin", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should declare all parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("source: close({"))
			Expect(cueOutput).To(ContainSubstring("hcl: string"))
			Expect(cueOutput).To(ContainSubstring(`remote: *"https://github.com/kubevela-contrib/terraform-modules.git" | string`))
			Expect(cueOutput).To(ContainSubstring("path?: string"))
			Expect(cueOutput).To(ContainSubstring("deleteResource: *true | bool"))
			Expect(cueOutput).To(ContainSubstring("forceDelete: *false | bool"))
			Expect(cueOutput).To(ContainSubstring("variable: {...}"))
			Expect(cueOutput).To(ContainSubstring("jobEnv?: {...}"))
			Expect(cueOutput).To(ContainSubstring("writeConnectionSecretToRef?: {"))
			Expect(cueOutput).To(ContainSubstring("providerRef?: {"))
			Expect(cueOutput).To(ContainSubstring("region?: string"))
		})

		It("should create a terraform Configuration resource with correct metadata and spec", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "terraform.core.oam.dev/v1beta2"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Configuration"`))
			Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepName)`))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring("deleteResource: parameter.deleteResource"))
			Expect(cueOutput).To(ContainSubstring("variable: parameter.variable"))
			Expect(cueOutput).To(ContainSubstring("forceDelete: parameter.forceDelete"))
			Expect(cueOutput).To(ContainSubstring("parameter.source.path != _|_"))
			Expect(cueOutput).To(ContainSubstring("path: parameter.source.path"))
			Expect(cueOutput).To(ContainSubstring("parameter.source.remote != _|_"))
			Expect(cueOutput).To(ContainSubstring("remote: parameter.source.remote"))
			Expect(cueOutput).To(ContainSubstring("parameter.source.hcl != _|_"))
			Expect(cueOutput).To(ContainSubstring("hcl: parameter.source.hcl"))
			Expect(cueOutput).To(ContainSubstring("parameter.providerRef != _|_"))
			Expect(cueOutput).To(ContainSubstring("providerRef: parameter.providerRef"))
			Expect(cueOutput).To(ContainSubstring("parameter.jobEnv != _|_"))
			Expect(cueOutput).To(ContainSubstring("jobEnv: parameter.jobEnv"))
			Expect(cueOutput).To(ContainSubstring("parameter.writeConnectionSecretToRef != _|_"))
			Expect(cueOutput).To(ContainSubstring("writeConnectionSecretToRef: parameter.writeConnectionSecretToRef"))
			Expect(cueOutput).To(ContainSubstring("parameter.region != _|_"))
			Expect(cueOutput).To(ContainSubstring("region: parameter.region"))
		})

		It("should apply the resource with kube.#Apply and wait for Available state", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("apply.$returns.value.status != _|_"))
			Expect(cueOutput).To(ContainSubstring("apply.$returns.value.status.apply != _|_"))
			Expect(cueOutput).To(ContainSubstring(`apply.$returns.value.status.apply.state == "Available"`))
		})

		It("should have exactly one kube.#Apply and one builtin.#ConditionalWait", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
		})
	})
})
