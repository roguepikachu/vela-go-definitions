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

var _ = Describe("CleanJobs WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.CleanJobs()
		Expect(step.GetName()).To(Equal("clean-jobs"))
		Expect(step.GetDescription()).To(Equal("clean applied jobs in the cluster"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.CleanJobs()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate the correct step header", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			Expect(cueOutput).To(ContainSubstring(`"clean-jobs": {`))
		})

		It("should import vela/kube", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
		})

		It("should declare the expected parameters", func() {
			Expect(cueOutput).To(ContainSubstring("labelselector?: {...}"))
			Expect(cueOutput).To(ContainSubstring("namespace: *context.namespace | string"))
		})

		It("should generate the cleanJobs kube.#Delete action", func() {
			Expect(cueOutput).To(ContainSubstring("cleanJobs: kube.#Delete & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Job"`))
			Expect(cueOutput).To(ContainSubstring("name: context.name"))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			Expect(cueOutput).To(ContainSubstring(`parameter["labelselector"] != _|_`))
			Expect(cueOutput).To(ContainSubstring("matchingLabels: parameter.labelselector"))
			Expect(cueOutput).To(ContainSubstring(`parameter["labelselector"] == _|_`))
			Expect(cueOutput).To(ContainSubstring(`"workflow.oam.dev/name": context.name`))
		})

		It("should generate the cleanPods kube.#Delete action", func() {
			Expect(cueOutput).To(ContainSubstring("cleanPods: kube.#Delete & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "pod"`))

			cleanPodsBlock := cueOutput[strings.Index(cueOutput, "cleanPods: kube.#Delete & {"):]
			Expect(cleanPodsBlock).To(ContainSubstring(`parameter["labelselector"] != _|_`))
			Expect(cleanPodsBlock).To(ContainSubstring("matchingLabels: parameter.labelselector"))
			Expect(cleanPodsBlock).To(ContainSubstring(`parameter["labelselector"] == _|_`))
			Expect(cleanPodsBlock).To(ContainSubstring(`"workflow.oam.dev/name": context.name`))
		})

		It("should be structurally correct with two delete actions and filters", func() {
			Expect(strings.Count(cueOutput, "kube.#Delete & {")).To(Equal(2))
			Expect(strings.Count(cueOutput, "filter: {")).To(Equal(2))
			Expect(strings.Count(cueOutput, `"workflow.oam.dev/name": context.name`)).To(Equal(2))
		})
	})
})
