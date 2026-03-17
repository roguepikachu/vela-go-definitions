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

var _ = Describe("PrintMessageInStatus WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.PrintMessageInStatus()
		Expect(step.GetName()).To(Equal("print-message-in-status"))
		Expect(step.GetDescription()).To(Equal("print message in workflow step status"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.PrintMessageInStatus()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Process Control"`))
		})

		It("should import vela/builtin", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should declare required message parameter", func() {
			Expect(cueOutput).To(ContainSubstring("message: string"))
		})

		It("should generate template with a single builtin.#Message call passing full parameter", func() {
			Expect(cueOutput).To(ContainSubstring("builtin.#Message & {"))
			Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			Expect(cueOutput).NotTo(MatchRegexp(`\$params:\s*\{`))
			Expect(strings.Count(cueOutput, "builtin.#Message & {")).To(Equal(1))
		})
	})
})
