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

var _ = Describe("Request WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.Request()
		Expect(step.GetName()).To(Equal("request"))
		Expect(step.GetDescription()).To(Equal("Send request to the url"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Request()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "External Integration"`))
			Expect(cueOutput).NotTo(ContainSubstring(`"alias":`))
		})

		It("should import vela/op, vela/http, and encoding/json", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/http"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
		})

		It("should declare all parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("url: string"))
			Expect(cueOutput).To(ContainSubstring(`method: *"GET" | "POST" | "PUT" | "DELETE"`))
			Expect(cueOutput).To(ContainSubstring("body?: {...}"))
			Expect(cueOutput).To(ContainSubstring("header?: [string]: string"))
		})

		It("should generate HTTP request template with conditional body, header, wait, and failure handling", func() {
			Expect(cueOutput).To(ContainSubstring("http.#HTTPDo & {"))
			Expect(cueOutput).To(ContainSubstring("method: parameter.method"))
			Expect(cueOutput).To(ContainSubstring("url: parameter.url"))
			Expect(cueOutput).To(ContainSubstring(`parameter["body"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("body: json.Marshal(parameter.body)"))
			Expect(cueOutput).To(ContainSubstring(`parameter["header"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("header: parameter.header"))
			Expect(cueOutput).To(ContainSubstring("op.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("continue: req.$returns != _|_"))
			Expect(cueOutput).To(ContainSubstring(`message?: "Waiting for response from \(parameter.url)"`))
			Expect(cueOutput).To(ContainSubstring("op.#Steps & {"))
			Expect(cueOutput).To(ContainSubstring("req.$returns.statusCode > 400"))
			Expect(cueOutput).To(ContainSubstring("requestFail: op.#Fail & {"))
			Expect(cueOutput).To(ContainSubstring(`message: "request of \(parameter.url) is fail: \(req.$returns.statusCode)"`))
			Expect(cueOutput).To(ContainSubstring("response: json.Unmarshal(req.$returns.body)"))
		})

		It("should have exactly one of each action type", func() {
			Expect(strings.Count(cueOutput, "http.#HTTPDo & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "op.#ConditionalWait & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "op.#Steps & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "op.#Fail & {")).To(Equal(1))
		})
	})
})
