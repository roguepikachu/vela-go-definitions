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
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ApplyDeployment()
			Expect(step.GetName()).To(Equal("apply-deployment"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ApplyDeployment()
			Expect(step.GetDescription()).To(Equal("Apply deployment with specified image and cmd."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyDeployment()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"apply-deployment": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required image", func() {
				Expect(cueOutput).To(ContainSubstring("image: string"))
			})

			It("should have replicas with default 1", func() {
				Expect(cueOutput).To(ContainSubstring("replicas: *1 | int"))
			})

			It("should have cluster with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			})

			It("should have optional cmd list", func() {
				Expect(cueOutput).To(ContainSubstring("cmd?: [...string]"))
			})
		})

		Describe("Template: Deployment resource", func() {
			It("should create apps/v1 Deployment", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "apps/v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Deployment"`))
			})

			It("should set metadata name from context.stepName", func() {
				Expect(cueOutput).To(ContainSubstring("name: context.stepName"))
			})

			It("should set metadata namespace from context", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should use step label for selector and pod labels", func() {
				count := strings.Count(cueOutput, `"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"`)
				Expect(count).To(Equal(2))
			})

			It("should set replicas from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("replicas: parameter.replicas"))
			})

			It("should set container name from stepName", func() {
				Expect(cueOutput).To(ContainSubstring("name: context.stepName"))
			})

			It("should set container image from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("image: parameter.image"))
			})

			It("should conditionally set command when cmd is set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["cmd"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("command: parameter.cmd"))
			})
		})

		Describe("Template: kube.#Apply", func() {
			It("should use kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should pass cluster parameter", func() {
				Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
			})
		})

		Describe("Template: wait action", func() {
			It("should use builtin.#ConditionalWait", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			})

			It("should wait for readyReplicas to match parameter", func() {
				Expect(cueOutput).To(ContainSubstring("output.$returns.value.status.readyReplicas == parameter.replicas"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
