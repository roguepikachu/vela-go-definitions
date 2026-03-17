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

var _ = Describe("ExportService WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.ExportService()
		Expect(step.GetName()).To(Equal("export-service"))
		Expect(step.GetDescription()).To(Equal("Export service to clusters specified by topology."))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ExportService()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, scope, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
			Expect(cueOutput).To(ContainSubstring(`"export-service": {`))
		})

		It("should import vela/op and vela/kube", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
		})

		It("should declare all parameters with correct types", func() {
			Expect(cueOutput).To(ContainSubstring("name?: string"))
			Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			Expect(cueOutput).To(ContainSubstring("ip: string"))
			Expect(cueOutput).To(ContainSubstring("port: int"))
			Expect(cueOutput).To(ContainSubstring("targetPort: int"))
			Expect(cueOutput).To(ContainSubstring("topology?: string"))
		})

		It("should define meta block with context defaults and conditional overrides", func() {
			Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
			Expect(cueOutput).To(ContainSubstring("*context.name | string"))
			Expect(cueOutput).To(ContainSubstring(`parameter["name"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
		})

		It("should define objects array with Service and Endpoints", func() {
			Expect(cueOutput).To(ContainSubstring("objects: ["))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Service"`))
			Expect(cueOutput).To(ContainSubstring("metadata: meta"))
			Expect(cueOutput).To(ContainSubstring(`type: "ClusterIP"`))
			Expect(cueOutput).To(ContainSubstring(`protocol: "TCP"`))
			Expect(cueOutput).To(ContainSubstring("port: parameter.port"))
			Expect(cueOutput).To(ContainSubstring("targetPort: parameter.targetPort"))
			Expect(cueOutput).To(ContainSubstring(`kind: "Endpoints"`))
			Expect(cueOutput).To(ContainSubstring("ip: parameter.ip"))
			Expect(cueOutput).To(ContainSubstring("port: parameter.targetPort"))
			count := strings.Count(cueOutput, "metadata: meta")
			Expect(count).To(Equal(2))
		})

		It("should get placements and apply via nested comprehension over objects", func() {
			Expect(cueOutput).To(ContainSubstring("op.#GetPlacementsFromTopologyPolicies & {"))
			Expect(cueOutput).To(ContainSubstring("policies: *[] | [...string]"))
			Expect(cueOutput).To(ContainSubstring("parameter.topology != _|_"))
			Expect(cueOutput).To(ContainSubstring("policies: [parameter.topology]"))
			Expect(cueOutput).To(ContainSubstring("for p in getPlacements.placements"))
			Expect(cueOutput).To(ContainSubstring("for o in objects"))
			Expect(cueOutput).To(ContainSubstring(`"\(p.cluster)-\(o.kind)"`))
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring("value:   o"))
			Expect(cueOutput).To(ContainSubstring("cluster: p.cluster"))
		})

		It("should be structurally correct", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "getPlacements:")).To(Equal(1))
		})
	})
})
