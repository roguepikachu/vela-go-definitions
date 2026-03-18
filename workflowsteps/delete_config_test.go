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

var _ = Describe("DeleteConfig WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.DeleteConfig()
		Expect(step.GetName()).To(Equal("delete-config"))
		Expect(step.GetDescription()).To(Equal("Delete a config"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.DeleteConfig()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate the correct step header", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Config Management"`))
			Expect(cueOutput).To(ContainSubstring(`"delete-config": {`))
		})

		It("should import vela/config", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
		})

		It("should declare all parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
		})

		It("should generate the template with a single config.#DeleteConfig call", func() {
			Expect(cueOutput).To(ContainSubstring("config.#DeleteConfig & {"))
			Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			count := strings.Count(cueOutput, "config.#DeleteConfig & {")
			Expect(count).To(Equal(1))
		})
	})
})
