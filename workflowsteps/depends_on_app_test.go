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

var _ = Describe("DependsOnApp WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.DependsOnApp()
		Expect(step.GetName()).To(Equal("depends-on-app"))
		Expect(step.GetDescription()).To(Equal("Wait for the specified Application to complete."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.DependsOnApp()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"depends-on-app": {`))
		})

		It("should import vela/kube, vela/builtin, and encoding/yaml", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/yaml"`))
		})

		It("should declare name and namespace as required string parameters", func() {
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("namespace: string"))
		})

		It("should read an Application resource via kube.#Read", func() {
			Expect(cueOutput).To(ContainSubstring("dependsOn: kube.#Read & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "core.oam.dev/v1beta1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Application"`))
			Expect(cueOutput).To(ContainSubstring("name:      parameter.name"))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
		})

		It("should define a load block with error path that reads configMap and applies unmarshaled template", func() {
			Expect(cueOutput).To(ContainSubstring("load: {"))
			Expect(cueOutput).To(ContainSubstring("dependsOn.$returns.err != _|_"))
			Expect(cueOutput).To(ContainSubstring("configMap: kube.#Read & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "ConfigMap"`))
			Expect(cueOutput).To(ContainSubstring(`configMap.$returns.value.data["application"]`))
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring("yaml.Unmarshal(template)"))
			Expect(cueOutput).To(ContainSubstring(`apply.$returns.value.status.status == "running"`))
		})

		It("should define a success path that waits for dependsOn status running", func() {
			Expect(cueOutput).To(ContainSubstring("dependsOn.$returns.err == _|_"))
			Expect(cueOutput).To(ContainSubstring(`dependsOn.$returns.value.status.status == "running"`))
		})

		It("should be structurally correct with expected action counts", func() {
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(2))
			Expect(strings.Count(cueOutput, "kube.#Read & {")).To(Equal(2))
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "dependsOn.$returns.err != _|_")).To(Equal(4))
			Expect(strings.Count(cueOutput, "dependsOn.$returns.err == _|_")).To(Equal(1))
		})
	})
})
