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

var _ = Describe("CollectServiceEndpoints WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.CollectServiceEndpoints()
		Expect(step.GetName()).To(Equal("collect-service-endpoints"))
		Expect(step.GetDescription()).To(Equal("Collect service endpoints for the application."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.CollectServiceEndpoints()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"collect-service-endpoints": {`))
		})

		It("should import vela/builtin, vela/query, and strconv", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/query"`))
			Expect(cueOutput).To(ContainSubstring(`"strconv"`))
		})

		It("should declare all parameter fields with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("name: *context.name | string"))
			Expect(cueOutput).To(ContainSubstring("namespace: *context.namespace | string"))
			Expect(cueOutput).To(ContainSubstring("components?: [...string]"))
			Expect(cueOutput).To(ContainSubstring("port?: int"))
			Expect(cueOutput).To(ContainSubstring("portName?: string"))
			Expect(cueOutput).To(ContainSubstring("outer?: bool"))
			Expect(cueOutput).To(ContainSubstring(`protocal: *"http" | "https"`))
		})

		It("should invoke query.#CollectServiceEndpoints with app params and conditional components filter", func() {
			Expect(cueOutput).To(ContainSubstring("query.#CollectServiceEndpoints & {"))
			Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			Expect(cueOutput).To(ContainSubstring(`parameter["components"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("components: parameter.components"))
		})

		It("should filter by portName when set, falling back to full list", func() {
			Expect(cueOutput).To(ContainSubstring("eps_port_name_filtered: *[] | [...]"))
			Expect(cueOutput).To(ContainSubstring(`parameter["portName"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("parameter.portName == ep.endpoint.portName"))
			Expect(cueOutput).To(ContainSubstring(`parameter["portName"] == _|_`))
			Expect(cueOutput).To(ContainSubstring("collect.$returns.list"))
		})

		It("should filter by port when set and alias result to eps", func() {
			Expect(cueOutput).To(ContainSubstring("eps_port_filtered: *[] | [...]"))
			Expect(cueOutput).To(ContainSubstring(`parameter["port"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("parameter.port == ep.endpoint.port"))
			Expect(cueOutput).To(ContainSubstring("eps: eps_port_filtered"))
		})

		It("should filter endpoints by outer flag when set, passing through otherwise", func() {
			Expect(cueOutput).To(ContainSubstring("endpoints: *[] | [...]"))
			Expect(cueOutput).To(ContainSubstring(`parameter["outer"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("tmps:"))
			Expect(cueOutput).To(ContainSubstring("ep.endpoint.inner == _|_"))
			Expect(cueOutput).To(ContainSubstring("outer: true"))
			Expect(cueOutput).To(ContainSubstring("outer: !ep.endpoint.inner"))
			Expect(cueOutput).To(ContainSubstring("!parameter.outer || ep.outer"))
			Expect(cueOutput).To(ContainSubstring(`parameter["outer"] == _|_`))
		})

		It("should use builtin.#ConditionalWait until endpoints are available", func() {
			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("len(outputs.endpoints) > 0"))
		})

		It("should extract first endpoint and build URL with protocal interpolation", func() {
			Expect(cueOutput).To(ContainSubstring("endpoint: outputs.endpoints[0].endpoint"))
			Expect(cueOutput).To(ContainSubstring("strconv.FormatInt(endpoint.port, 10)"))
			Expect(cueOutput).To(ContainSubstring(`\(parameter.protocal)`))
			Expect(cueOutput).To(ContainSubstring(`\(endpoint.host)`))
			Expect(cueOutput).To(ContainSubstring(`\(_portStr)`))
		})

		It("should be structurally correct with one collect and one wait action", func() {
			Expect(strings.Count(cueOutput, "query.#CollectServiceEndpoints & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))

			lines := strings.Split(cueOutput, "\n")
			inApp := false
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "app: {" {
					inApp = true
				}
				if inApp && trimmed == "}" {
					break
				}
				if inApp {
					Expect(trimmed).NotTo(ContainSubstring(`parameter["name"] != _|_`))
					Expect(trimmed).NotTo(ContainSubstring(`parameter["name"] == _|_`))
					Expect(trimmed).NotTo(ContainSubstring(`parameter["namespace"] != _|_`))
					Expect(trimmed).NotTo(ContainSubstring(`parameter["namespace"] == _|_`))
				}
			}
		})
	})
})
