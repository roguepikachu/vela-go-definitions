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

var _ = Describe("ReadObject WorkflowStep", func() {
	It("should have correct name and description", func() {
		step := workflowsteps.ReadObject()
		Expect(step.GetName()).To(Equal("read-object"))
		Expect(step.GetDescription()).To(Equal("Read Kubernetes objects from cluster for your workflow steps"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ReadObject()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
		})

		It("should import vela/kube", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
		})

		It("should declare all parameters with correct types, defaults, and descriptions", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: *"core.oam.dev/v1beta1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: *"Application"`))
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring(`namespace: *"default"`))
			Expect(cueOutput).To(ContainSubstring(`cluster: *""`))

			Expect(cueOutput).To(ContainSubstring("apiVersion of the object"))
			Expect(cueOutput).To(ContainSubstring("name of the object"))
			Expect(cueOutput).To(ContainSubstring("namespace of the resource"))
			Expect(cueOutput).To(ContainSubstring("cluster you want to apply"))
		})

		It("should generate a single kube.#Read template that passes all parameters", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Read & {"))
			Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
			Expect(cueOutput).To(ContainSubstring("apiVersion: parameter.apiVersion"))
			Expect(cueOutput).To(ContainSubstring("kind: parameter.kind"))
			Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))

			count := strings.Count(cueOutput, "kube.#Read & {")
			Expect(count).To(Equal(1))
		})
	})
})
