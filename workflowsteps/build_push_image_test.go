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

var _ = Describe("BuildPushImage WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.BuildPushImage()
		Expect(step.GetName()).To(Equal("build-push-image"))
		Expect(step.GetDescription()).To(Equal("Build and push image from git url"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.BuildPushImage()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, quoted name, and empty alias", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "CI Integration"`))
			Expect(cueOutput).To(ContainSubstring(`"build-push-image": {`))
			Expect(cueOutput).To(ContainSubstring(`alias: ""`))
		})

		It("should import all required packages", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			Expect(cueOutput).To(ContainSubstring(`"strings"`))
		})

		It("should define #secret and #git helper types", func() {
			Expect(cueOutput).To(ContainSubstring("#secret: {"))
			secretIdx := strings.Index(cueOutput, "#secret: {")
			secretBlock := cueOutput[secretIdx : secretIdx+60]
			Expect(secretBlock).To(ContainSubstring("name: string"))
			Expect(secretBlock).To(ContainSubstring("key: string"))

			Expect(cueOutput).To(ContainSubstring("#git: {"))
			Expect(cueOutput).To(ContainSubstring("git: string"))
			Expect(cueOutput).To(ContainSubstring(`branch: *"master" | string`))
		})

		It("should declare all parameters with correct types, defaults, and credentials structure", func() {
			Expect(cueOutput).To(ContainSubstring(`kanikoExecutor: *"oamdev/kaniko-executor:v1.9.1" | string`))
			Expect(cueOutput).To(ContainSubstring(`dockerfile: *"./Dockerfile" | string`))
			Expect(cueOutput).To(ContainSubstring("image: string"))
			Expect(cueOutput).To(ContainSubstring("platform?: string"))
			Expect(cueOutput).To(ContainSubstring("buildArgs?: [...string]"))
			Expect(cueOutput).To(ContainSubstring(`verbosity: *"info" | "panic" | "fatal" | "error" | "warn" | "debug" | "trace"`))
			Expect(cueOutput).To(ContainSubstring("context: #git | string"))

			credIdx := strings.Index(cueOutput, "credentials?: {")
			Expect(credIdx).To(BeNumerically(">", 0))
			credBlock := cueOutput[credIdx:]
			Expect(credBlock).To(ContainSubstring("git?: {"))
			Expect(credBlock).To(ContainSubstring("image?: {"))
			imageIdx := strings.Index(credBlock, "image?: {")
			imageBlock := credBlock[imageIdx:]
			Expect(imageBlock).To(ContainSubstring(`key: *".dockerconfigjson" | string`))
		})

		It("should handle git context URL building with branch and prefix trimming", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.context.git != _|_"))
			Expect(cueOutput).To(ContainSubstring("parameter.context.git == _|_"))
			Expect(cueOutput).To(ContainSubstring(`strings.TrimPrefix(parameter.context.git, "git://")`))
			Expect(cueOutput).To(ContainSubstring(`git://\(address)#refs/heads/\(parameter.context.branch)`))
			Expect(cueOutput).To(ContainSubstring("value: parameter.context"))
		})

		It("should build kaniko Pod with correct spec, container args, and conditional mounts", func() {
			Expect(cueOutput).To(ContainSubstring("kaniko: kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Pod"`))
			Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepSessionID)-kaniko`))
			Expect(cueOutput).To(ContainSubstring(`--dockerfile=\(parameter.dockerfile)`))
			Expect(cueOutput).To(ContainSubstring(`--context=\(url.value)`))
			Expect(cueOutput).To(ContainSubstring(`--destination=\(parameter.image)`))
			Expect(cueOutput).To(ContainSubstring(`--verbosity=\(parameter.verbosity)`))
			Expect(cueOutput).To(ContainSubstring("parameter.platform != _|_"))
			Expect(cueOutput).To(ContainSubstring(`--customPlatform=\(parameter.platform)`))
			Expect(cueOutput).To(ContainSubstring("parameter.buildArgs != _|_"))
			Expect(cueOutput).To(ContainSubstring(`--build-arg=\(arg)`))
			Expect(cueOutput).To(ContainSubstring("image: parameter.kanikoExecutor"))
			Expect(cueOutput).To(ContainSubstring("parameter.credentials.image != _|_"))
			Expect(cueOutput).To(ContainSubstring(`mountPath: "/kaniko/.docker/"`))
			Expect(cueOutput).To(ContainSubstring("parameter.credentials.git != _|_"))
			Expect(cueOutput).To(ContainSubstring(`name: "GIT_TOKEN"`))
			Expect(cueOutput).To(ContainSubstring("secretKeyRef:"))
			Expect(cueOutput).To(ContainSubstring(`restartPolicy: "Never"`))
		})

		It("should have log, read, and wait actions for kaniko Pod lifecycle", func() {
			Expect(cueOutput).To(ContainSubstring("util.#Log & {"))
			Expect(cueOutput).To(ContainSubstring("read: kube.#Read & {"))

			// Verify read targets a Pod (comment 5: prevent read-kind regressions)
			readIdx := strings.Index(cueOutput, "read: kube.#Read & {")
			readBlock := cueOutput[readIdx : readIdx+200]
			Expect(readBlock).To(ContainSubstring(`kind: "Pod"`))

			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("read.$returns.value.status != _|_"))
			Expect(cueOutput).To(ContainSubstring(`read.$returns.value.status.phase == "Succeeded"`))

			count := strings.Count(cueOutput, `\(context.name)-\(context.stepSessionID)-kaniko`)
			Expect(count).To(BeNumerically(">=", 3))
		})

		It("should have exactly one of each action type", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "kube.#Read & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "util.#Log & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
		})
	})
})
