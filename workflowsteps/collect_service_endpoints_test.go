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
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.CollectServiceEndpoints()
			Expect(step.GetName()).To(Equal("collect-service-endpoints"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.CollectServiceEndpoints()
			Expect(step.GetDescription()).To(Equal("Collect service endpoints for the application."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.CollectServiceEndpoints()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"collect-service-endpoints": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})

			It("should import vela/query", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/query"`))
			})

			It("should import strconv", func() {
				Expect(cueOutput).To(ContainSubstring(`"strconv"`))
			})
		})

		Describe("Parameters", func() {
			It("should have name with context.name default", func() {
				Expect(cueOutput).To(ContainSubstring("name: *context.name | string"))
			})

			It("should have namespace with context.namespace default", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: *context.namespace | string"))
			})

			It("should have optional components list", func() {
				Expect(cueOutput).To(ContainSubstring("components?: [...string]"))
			})

			It("should have optional port", func() {
				Expect(cueOutput).To(ContainSubstring("port?: int"))
			})

			It("should have optional portName", func() {
				Expect(cueOutput).To(ContainSubstring("portName?: string"))
			})

			It("should have optional outer", func() {
				Expect(cueOutput).To(ContainSubstring("outer?: bool"))
			})

			It("should have protocal with http default", func() {
				Expect(cueOutput).To(ContainSubstring(`protocal: *"http" | "https"`))
			})
		})

		Describe("Template: collect action", func() {
			It("should use query.#CollectServiceEndpoints", func() {
				Expect(cueOutput).To(ContainSubstring("query.#CollectServiceEndpoints & {"))
			})

			It("should pass name directly to app", func() {
				Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			})

			It("should pass namespace directly to app", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})

			It("should conditionally pass components in filter", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["components"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("components: parameter.components"))
			})
		})

		Describe("Template: outputs filtering", func() {
			It("should define eps_port_name_filtered with empty default", func() {
				Expect(cueOutput).To(ContainSubstring("eps_port_name_filtered: *[] | [...]"))
			})

			It("should filter by portName when set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["portName"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("parameter.portName == ep.endpoint.portName"))
			})

			It("should use collect.$returns.list when portName not set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["portName"] == _|_`))
				Expect(cueOutput).To(ContainSubstring("collect.$returns.list"))
			})

			It("should define eps_port_filtered with empty default", func() {
				Expect(cueOutput).To(ContainSubstring("eps_port_filtered: *[] | [...]"))
			})

			It("should filter by port when set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["port"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("parameter.port == ep.endpoint.port"))
			})

			It("should alias eps to eps_port_filtered", func() {
				Expect(cueOutput).To(ContainSubstring("eps: eps_port_filtered"))
			})

			It("should define endpoints with empty default", func() {
				Expect(cueOutput).To(ContainSubstring("endpoints: *[] | [...]"))
			})

			It("should compute tmps with outer flag when outer is set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["outer"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("tmps:"))
				Expect(cueOutput).To(ContainSubstring("ep.endpoint.inner == _|_"))
				Expect(cueOutput).To(ContainSubstring("outer: true"))
				Expect(cueOutput).To(ContainSubstring("outer: !ep.endpoint.inner"))
			})

			It("should filter endpoints by outer flag", func() {
				Expect(cueOutput).To(ContainSubstring("!parameter.outer || ep.outer"))
			})

			It("should use eps_port_filtered when outer not set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["outer"] == _|_`))
			})
		})

		Describe("Template: wait action", func() {
			It("should use builtin.#ConditionalWait", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			})

			It("should wait for endpoints length > 0", func() {
				Expect(cueOutput).To(ContainSubstring("len(outputs.endpoints) > 0"))
			})
		})

		Describe("Template: value block", func() {
			It("should extract first endpoint", func() {
				Expect(cueOutput).To(ContainSubstring("endpoint: outputs.endpoints[0].endpoint"))
			})

			It("should format port as string", func() {
				Expect(cueOutput).To(ContainSubstring("strconv.FormatInt(endpoint.port, 10)"))
			})

			It("should build URL with protocal interpolation", func() {
				Expect(cueOutput).To(ContainSubstring(`\(parameter.protocal)`))
				Expect(cueOutput).To(ContainSubstring(`\(endpoint.host)`))
				Expect(cueOutput).To(ContainSubstring(`\(_portStr)`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one query.#CollectServiceEndpoints", func() {
				count := strings.Count(cueOutput, "query.#CollectServiceEndpoints & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})

			It("should NOT have conditional name/namespace in app params", func() {
				// name and namespace have defaults, so no guards needed in app
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
})
