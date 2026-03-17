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

var _ = Describe("CheckMetrics WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.CheckMetrics()
		Expect(step.GetName()).To(Equal("check-metrics"))
		Expect(step.GetDescription()).To(Equal("Verify application's metrics"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.CheckMetrics()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate the correct step header", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"catalog": "Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"check-metrics": {`))
		})

		It("should import required packages", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/metrics"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should declare all parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("query: string"))
			Expect(cueOutput).To(ContainSubstring(`metricEndpoint?: "http://prometheus-server.o11y-system.svc:9090" | string`))
			Expect(cueOutput).To(ContainSubstring("condition: string"))
			Expect(cueOutput).To(ContainSubstring(`duration?: *"5m" | string`))
			Expect(cueOutput).To(ContainSubstring(`failDuration?: *"2m" | string`))
		})

		It("should generate the check action using metrics.#PromCheck", func() {
			Expect(cueOutput).To(ContainSubstring("metrics.#PromCheck & {"))
			Expect(cueOutput).To(ContainSubstring("query: parameter.query"))
			Expect(cueOutput).To(ContainSubstring("metricEndpoint: parameter.metricEndpoint"))
			Expect(cueOutput).To(ContainSubstring("condition: parameter.condition"))
			Expect(cueOutput).To(ContainSubstring("duration: parameter.duration"))
			Expect(cueOutput).To(ContainSubstring("failDuration: parameter.failDuration"))
		})

		It("should generate the fail block with guarded breakWorkflow", func() {
			Expect(cueOutput).To(ContainSubstring("check.$returns.failed != _|_"))
			Expect(cueOutput).To(ContainSubstring("check.$returns.failed == true"))
			Expect(cueOutput).To(ContainSubstring("builtin.#Fail & {"))
			Expect(cueOutput).To(ContainSubstring("message: check.$returns.message"))
		})

		It("should generate the wait action using builtin.#ConditionalWait", func() {
			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("continue: check.$returns.result"))
			Expect(cueOutput).To(ContainSubstring("check.$returns.message != _|_"))
			Expect(cueOutput).To(ContainSubstring("message: check.$returns.message"))
		})

		It("should have exactly one instance of each action type", func() {
			Expect(strings.Count(cueOutput, "metrics.#PromCheck & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#Fail & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
		})
	})
})
