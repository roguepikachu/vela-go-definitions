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

var _ = Describe("ShareCloudResource WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.ShareCloudResource()
		Expect(step.GetName()).To(Equal("share-cloud-resource"))
		Expect(step.GetDescription()).To(ContainSubstring("Sync secrets created by terraform component"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ShareCloudResource()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and scope", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
		})

		It("should import vela/op", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
		})

		It("should declare placements, policy, and env parameters", func() {
			Expect(cueOutput).To(ContainSubstring("placements: [...{"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Declare the location to bind"))

			placementsIdx := strings.Index(cueOutput, "placements: [...{")
			Expect(placementsIdx).To(BeNumerically(">", 0))
			placementsBlock := cueOutput[placementsIdx:]
			Expect(placementsBlock).To(ContainSubstring("namespace?: string"))
			Expect(placementsBlock).To(ContainSubstring("cluster?:"))

			Expect(cueOutput).To(ContainSubstring(`policy: *"" | string`))
			Expect(cueOutput).To(ContainSubstring("env: string"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Declare the name of the env in policy"))
		})

		It("should invoke a single op.#ShareCloudResource with direct field bindings", func() {
			Expect(cueOutput).To(ContainSubstring("app: op.#ShareCloudResource & {"))
			Expect(cueOutput).To(ContainSubstring("env: parameter.env"))
			Expect(cueOutput).To(ContainSubstring("policy: parameter.policy"))
			Expect(cueOutput).To(ContainSubstring("placements: parameter.placements"))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring("name: context.name"))
			Expect(cueOutput).NotTo(ContainSubstring("$params:"))
			Expect(strings.Count(cueOutput, "op.#ShareCloudResource & {")).To(Equal(1))
		})
	})
})
