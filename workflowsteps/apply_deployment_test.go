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

var _ = Describe("ApplyDeployment WorkflowStep", func() {
	It("should have correct name and description", func() {
		step := workflowsteps.ApplyDeployment()
		Expect(step.GetName()).To(Equal("apply-deployment"))
		Expect(step.GetDescription()).To(Equal("Apply deployment with specified image and cmd."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyDeployment()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			Expect(cueOutput).To(ContainSubstring(`"apply-deployment": {`))
		})

		It("should import vela/kube and vela/builtin", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should declare image, replicas, cluster, and cmd parameters with correct types and defaults", func() {
			Expect(cueOutput).To(ContainSubstring("image: string"))
			Expect(cueOutput).To(ContainSubstring("replicas: *1 | int"))
			Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			Expect(cueOutput).To(ContainSubstring("cmd?: [...string]"))
		})

		It("should build the Deployment resource with correct metadata, selector, pod spec, and conditional cmd", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "apps/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Deployment"`))
			Expect(cueOutput).To(ContainSubstring("name: context.stepName"))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring("replicas: parameter.replicas"))
			Expect(cueOutput).To(ContainSubstring("image: parameter.image"))
			Expect(cueOutput).To(ContainSubstring(`parameter["cmd"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("command: parameter.cmd"))

			count := strings.Count(cueOutput, `"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"`)
			Expect(count).To(Equal(2))
		})

		It("should apply the resource via kube.#Apply with cluster parameter", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
		})

		It("should wait for readyReplicas via builtin.#ConditionalWait", func() {
			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("output.$returns.value.status.readyReplicas == parameter.replicas"))
		})

		It("should have exactly one kube.#Apply and one builtin.#ConditionalWait", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
		})
	})
})
