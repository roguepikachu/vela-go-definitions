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

var _ = Describe("Webhook WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.Webhook()
		Expect(step.GetName()).To(Equal("webhook"))
		Expect(step.GetDescription()).To(ContainSubstring("Send a POST request to the specified Webhook URL"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Webhook()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "External Intergration"`))
		})

		It("should import all required packages", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/http"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
		})

		It("should declare url as ClosedUnion with value and secretRef options", func() {
			Expect(cueOutput).To(ContainSubstring("url: close({"))
			Expect(cueOutput).To(ContainSubstring("value: string"))
			Expect(cueOutput).To(ContainSubstring("}) | close({"))
			Expect(cueOutput).To(ContainSubstring("secretRef: {"))
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("key: string"))
			Expect(cueOutput).To(ContainSubstring("// +usage=name is the name of the secret"))
			Expect(cueOutput).To(ContainSubstring("// +usage=key is the key in the secret"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Specify the webhook url"))
		})

		It("should declare optional data parameter", func() {
			Expect(cueOutput).To(ContainSubstring("data?: {...}"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Specify the data you want to send"))
		})

		It("should read Application when no data provided and marshal accordingly", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Read & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "core.oam.dev/v1beta1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Application"`))
			Expect(cueOutput).To(ContainSubstring("name:      context.name"))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring("json.Marshal(read.$returns.value)"))
			Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.data)"))
			Expect(cueOutput).To(ContainSubstring("parameter.data == _|_"))
			Expect(cueOutput).To(ContainSubstring("parameter.data != _|_"))
		})

		It("should POST via http.#HTTPDo when url.value is set", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.url.value != _|_"))
			Expect(cueOutput).To(ContainSubstring("http.#HTTPDo & {"))
			Expect(cueOutput).To(ContainSubstring(`method: "POST"`))
			Expect(cueOutput).To(ContainSubstring("url:    parameter.url.value"))
			Expect(cueOutput).To(ContainSubstring(`header: "Content-Type": "application/json"`))
			Expect(cueOutput).To(ContainSubstring("body: data.value"))
		})

		It("should read Secret and POST via resolved URL when secretRef is set", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.url.secretRef != _|_ && parameter.url.value == _|_"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Secret"`))
			Expect(cueOutput).To(ContainSubstring("name:      parameter.url.secretRef.name"))
			Expect(cueOutput).To(ContainSubstring("util.#ConvertString & {"))
			Expect(cueOutput).To(ContainSubstring("base64.Decode(null, read.$returns.value.data[parameter.url.secretRef.key])"))
			Expect(cueOutput).To(ContainSubstring("url:    stringValue.$returns.str"))
		})

		It("should have correct structural counts", func() {
			Expect(strings.Count(cueOutput, "http.#HTTPDo & {")).To(Equal(2))
			Expect(strings.Count(cueOutput, "kube.#Read & {")).To(Equal(2))
			Expect(strings.Count(cueOutput, "util.#ConvertString & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "json.Marshal(")).To(Equal(2))
		})
	})
})
