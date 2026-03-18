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

var _ = Describe("RestartWorkflow WorkflowStep", func() {
	It("should have the correct name, description, category, and scope", func() {
		step := workflowsteps.RestartWorkflow()
		Expect(step.GetName()).To(Equal("restart-workflow"))
		Expect(step.GetDescription()).To(Equal("Schedule the current Application's workflow to restart at a specific time, after a duration, or at recurring intervals"))
		cue := step.ToCue()
		Expect(cue).To(ContainSubstring(`"category": "Workflow Control"`))
		Expect(cue).To(ContainSubstring(`"scope": "Application"`))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.RestartWorkflow()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header and imports", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should declare at, after, and every as optional string parameters with descriptions", func() {
			Expect(cueOutput).To(ContainSubstring("at?: string"))
			Expect(cueOutput).To(ContainSubstring("after?: string"))
			Expect(cueOutput).To(ContainSubstring("every?: string"))
			Expect(cueOutput).To(ContainSubstring("Schedule restart at a specific RFC3339 timestamp"))
			Expect(cueOutput).To(ContainSubstring("Schedule restart after a relative duration"))
			Expect(cueOutput).To(ContainSubstring("Schedule recurring restarts"))
		})

		It("should validate exactly one parameter and set _script conditionally for each", func() {
			Expect(cueOutput).To(ContainSubstring("_paramCount"))
			Expect(cueOutput).To(ContainSubstring("if parameter.at != _|_"))
			Expect(cueOutput).To(ContainSubstring("if parameter.after != _|_"))
			Expect(cueOutput).To(ContainSubstring("if parameter.every != _|_"))
			Expect(cueOutput).To(ContainSubstring("_script: string"))
			Expect(cueOutput).To(ContainSubstring("app.oam.dev/restart-workflow"))
			Expect(cueOutput).To(ContainSubstring("Convert duration to seconds"))
			Expect(cueOutput).To(ContainSubstring("builtin.#Fail"))
			Expect(cueOutput).To(ContainSubstring("Exactly one of"))
			Expect(cueOutput).To(ContainSubstring("_paramCount != 1"))
		})

		It("should create a Job via kube.#Apply with kubectl annotate container", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Job"`))
			Expect(cueOutput).To(ContainSubstring("restart-workflow"))
			Expect(cueOutput).To(ContainSubstring("context.stepSessionID"))
			Expect(cueOutput).To(ContainSubstring(`image:   "bitnami/kubectl:latest"`))
			Expect(cueOutput).To(ContainSubstring(`name:    "kubectl-annotate"`))
			Expect(cueOutput).To(ContainSubstring(`serviceAccountName: "kubevela-vela-core"`))
		})

		It("should wait for job completion with status checks", func() {
			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("continue: job.$returns.value.status.succeeded > 0"))
			Expect(cueOutput).To(ContainSubstring("job.$returns.value.status != _|_"))
			Expect(cueOutput).To(ContainSubstring("job.$returns.value.status.succeeded != _|_"))
		})

		It("should have exactly one of each action type", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#Fail & {")).To(Equal(1))
		})
	})
})
