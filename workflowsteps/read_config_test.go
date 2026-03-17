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

var _ = Describe("ReadConfig WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.ReadConfig()
		Expect(step.GetName()).To(Equal("read-config"))
		Expect(step.GetDescription()).To(Equal("Read a config"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ReadConfig()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Config Management"`))
		})

		It("should import vela/config", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
		})

		It("should declare name and namespace parameters with descriptions", func() {
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
			Expect(cueOutput).To(ContainSubstring("name of the config"))
			Expect(cueOutput).To(ContainSubstring("namespace of the config"))
		})

		It("should generate template with a single config.#ReadConfig call passing full parameter", func() {
			Expect(cueOutput).To(ContainSubstring("config.#ReadConfig & {"))
			Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			Expect(cueOutput).NotTo(MatchRegexp(`\$params:\s*\{`))
			Expect(strings.Count(cueOutput, "config.#ReadConfig & {")).To(Equal(1))
		})
	})
})
