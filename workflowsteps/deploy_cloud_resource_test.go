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

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"
	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("DeployCloudResource WorkflowStep", func() {
	var step *defkit.WorkflowStepDefinition

	BeforeEach(func() {
		step = workflowsteps.DeployCloudResource()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(step.GetName()).To(Equal("deploy-cloud-resource"))
		})

		It("should have the correct description", func() {
			Expect(step.GetDescription()).To(Equal(
				"Deploy cloud resource and deliver secret to multi clusters.",
			))
		})

		It("should be a workflow step definition type", func() {
			Expect(step.DefType()).To(Equal(defkit.DefinitionTypeWorkflowStep))
		})

		It("should have DefName matching GetName", func() {
			Expect(step.DefName()).To(Equal(step.GetName()))
		})

		It("should have the correct category", func() {
			Expect(step.GetCategory()).To(Equal("Application Delivery"))
		})

		It("should have the correct scope", func() {
			Expect(step.GetScope()).To(Equal("Application"))
		})

		It("should import vela/op", func() {
			Expect(step.GetImports()).To(ConsistOf("vela/op"))
		})
	})

	Describe("Parameters", func() {
		It("should have exactly 2 top-level parameters", func() {
			Expect(step.GetParams()).To(HaveLen(2))
		})

		It("should have parameters in correct order", func() {
			params := step.GetParams()
			Expect(params[0].Name()).To(Equal("policy"))
			Expect(params[1].Name()).To(Equal("env"))
		})

		Describe("policy parameter", func() {
			It("should have a default of empty string", func() {
				Expect(step.GetParams()[0]).NotTo(BeRequired())
			})

			It("should be a StringParam", func() {
				_, ok := step.GetParams()[0].(*defkit.StringParam)
				Expect(ok).To(BeTrue(), "policy should be a StringParam")
			})

			It("should have correct description", func() {
				Expect(step.GetParams()[0]).To(HaveDescription(
					"Declare the name of the env-binding policy, if empty, the first env-binding policy will be used",
				))
			})
		})

		Describe("env parameter", func() {
			It("should be required", func() {
				Expect(step.GetParams()[1]).To(BeRequired())
			})

			It("should be a StringParam", func() {
				_, ok := step.GetParams()[1].(*defkit.StringParam)
				Expect(ok).To(BeTrue(), "env should be a StringParam")
			})

			It("should have correct description", func() {
				Expect(step.GetParams()[1]).To(HaveDescription(
					"Declare the name of the env in policy",
				))
			})
		})
	})

	Describe("Helper Definitions", func() {
		It("should have no helper definitions", func() {
			Expect(step.GetHelperDefinitions()).To(BeEmpty())
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Import block", func() {
			It("should import vela/op", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			})
		})

		Describe("Step header", func() {
			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"deploy-cloud-resource": {`))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should have correct category annotation", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			})

			It("should have correct scope label", func() {
				Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					`description: "Deploy cloud resource and deliver secret to multi clusters."`,
				))
			})
		})

		Describe("Template block", func() {
			It("should have a template section", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should use op.#DeployCloudResource builtin", func() {
				Expect(cueOutput).To(ContainSubstring("app: op.#DeployCloudResource & {"))
			})

			It("should pass env parameter directly (no $params wrapper)", func() {
				Expect(cueOutput).To(ContainSubstring("env: parameter.env"))
				Expect(cueOutput).NotTo(ContainSubstring("$params"))
			})

			It("should pass policy parameter directly", func() {
				Expect(cueOutput).To(ContainSubstring("policy: parameter.policy"))
			})

			It("should pass context.namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should pass context.name", func() {
				Expect(cueOutput).To(ContainSubstring("name: context.name"))
			})
		})

		Describe("Parameter block", func() {
			It("should have a parameter section", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: {"))
			})

			It("should have policy with default empty string", func() {
				Expect(cueOutput).To(ContainSubstring(`policy: *"" | string`))
			})

			It("should have env as required string", func() {
				Expect(cueOutput).To(ContainSubstring("env: string"))
			})

			It("should include usage comment for policy", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Declare the name of the env-binding policy"))
			})

			It("should include usage comment for env", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Declare the name of the env in policy"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, `"deploy-cloud-resource": {`)
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})

			It("should have import before header", func() {
				importIdx := strings.Index(cueOutput, "import (")
				headerIdx := strings.Index(cueOutput, `"deploy-cloud-resource": {`)
				Expect(importIdx).To(BeNumerically("<", headerIdx))
			})
		})
	})

	Describe("YAML Generation", func() {
		It("should produce valid YAML with correct structure", func() {
			yamlBytes, err := step.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			yamlStr := string(yamlBytes)

			Expect(yamlStr).To(ContainSubstring("apiVersion: core.oam.dev/v1beta1"))
			Expect(yamlStr).To(ContainSubstring("kind: WorkflowStepDefinition"))
			Expect(yamlStr).To(ContainSubstring("name: deploy-cloud-resource"))
		})
	})
})
