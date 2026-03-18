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

var _ = Describe("DeployCloudResource WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.DeployCloudResource()
		Expect(step.GetName()).To(Equal("deploy-cloud-resource"))
		Expect(step.GetDescription()).To(Equal("Deploy cloud resource and deliver secret to multi clusters."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.DeployCloudResource()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct type, category, scope, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
			Expect(cueOutput).To(ContainSubstring(`"deploy-cloud-resource": {`))
		})

		It("should import vela/op", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
		})

		It("should declare policy with empty default and required env", func() {
			Expect(cueOutput).To(ContainSubstring(`policy: *"" | string`))
			Expect(cueOutput).To(ContainSubstring("env: string"))
		})

		It("should invoke op.#DeployCloudResource with correct field bindings", func() {
			Expect(cueOutput).To(ContainSubstring("app: op.#DeployCloudResource & {"))
			Expect(cueOutput).To(ContainSubstring("env: parameter.env"))
			Expect(cueOutput).To(ContainSubstring("policy: parameter.policy"))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring("name: context.name"))
			Expect(cueOutput).NotTo(ContainSubstring("$params"))
		})

		It("should have exactly one op.#DeployCloudResource invocation", func() {
			count := strings.Count(cueOutput, "op.#DeployCloudResource & {")
			Expect(count).To(Equal(1))
		})
	})
})
