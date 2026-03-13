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
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.DeleteConfig()
			Expect(step.GetName()).To(Equal("delete-config"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.DeleteConfig()
			Expect(step.GetDescription()).To(Equal("Delete a config"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.DeleteConfig()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Config Management"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"delete-config": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/config", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required name", func() {
				Expect(cueOutput).To(ContainSubstring("name: string"))
			})

			It("should have namespace with default context.namespace", func() {
				Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
			})
		})

		Describe("Template", func() {
			It("should call config.#DeleteConfig", func() {
				Expect(cueOutput).To(ContainSubstring("config.#DeleteConfig & {"))
			})

			It("should pass name parameter", func() {
				Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			})

			It("should pass namespace parameter", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one config.#DeleteConfig", func() {
				count := strings.Count(cueOutput, "config.#DeleteConfig & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
