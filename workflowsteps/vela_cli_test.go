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

var _ = Describe("VelaCli WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.VelaCli()
		Expect(step.GetName()).To(Equal("vela-cli"))
		Expect(step.GetDescription()).To(Equal("Run a vela command"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.VelaCli()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Scripts & Commands"`))
		})

		It("should import vela/kube, vela/builtin, and vela/util", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
		})

		It("should declare all parameters with correct types, defaults, and descriptions", func() {
			Expect(cueOutput).To(ContainSubstring("command: [...string]"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Specify the vela command"))
			Expect(cueOutput).To(ContainSubstring(`image: *"oamdev/vela-cli:v1.6.4" | string`))
			Expect(cueOutput).To(ContainSubstring("// +usage=Specify the image"))
			Expect(cueOutput).To(ContainSubstring(`serviceAccountName: *"kubevela-vela-core" | string`))
			Expect(cueOutput).To(ContainSubstring("// +usage=specify serviceAccountName want to use"))
			Expect(cueOutput).To(ContainSubstring("storage?:"))
			Expect(cueOutput).To(ContainSubstring("secret?:"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Mount Secret type storage"))
			Expect(cueOutput).To(ContainSubstring("hostPath?:"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Declare host path type storage"))
			Expect(cueOutput).To(ContainSubstring("secretName: string"))
			Expect(cueOutput).To(ContainSubstring("defaultMode: *420 | int"))
			Expect(cueOutput).To(ContainSubstring(`*"Directory"`))
			Expect(cueOutput).To(ContainSubstring(`"DirectoryOrCreate"`))
			Expect(cueOutput).To(ContainSubstring(`"FileOrCreate"`))
			Expect(cueOutput).To(ContainSubstring("items?:"))
			Expect(cueOutput).To(ContainSubstring("mode: *511 | int"))
		})

		It("should generate mountsArray with guarded for-each over secret and hostPath storage", func() {
			Expect(cueOutput).To(ContainSubstring("mountsArray:"))
			Expect(cueOutput).To(ContainSubstring("parameter.storage != _|_ && parameter.storage.secret != _|_"))
			Expect(cueOutput).To(ContainSubstring(`"secret-" + m.name`))
			Expect(cueOutput).To(ContainSubstring("mountPath: m.mountPath"))
			Expect(cueOutput).To(ContainSubstring("if m.subPath != _|_"))
			Expect(cueOutput).To(ContainSubstring("subPath: m.subPath"))
			Expect(cueOutput).To(ContainSubstring("parameter.storage != _|_ && parameter.storage.hostPath != _|_"))
			Expect(cueOutput).To(ContainSubstring(`"hostpath-" + m.name`))
		})

		It("should generate volumesList with secret and hostPath volumes", func() {
			Expect(cueOutput).To(ContainSubstring("volumesList:"))
			Expect(cueOutput).To(ContainSubstring("defaultMode: m.defaultMode"))
			Expect(cueOutput).To(ContainSubstring("secretName: m.secretName"))
			Expect(cueOutput).To(ContainSubstring("if m.items != _|_"))
			Expect(cueOutput).To(ContainSubstring("items: m.items"))
			Expect(cueOutput).To(ContainSubstring("path: m.path"))
		})

		It("should generate deDupVolumesArray with dedup pattern", func() {
			Expect(cueOutput).To(ContainSubstring("deDupVolumesArray:"))
			Expect(cueOutput).To(ContainSubstring("for val in ["))
			Expect(cueOutput).To(ContainSubstring("for i, vi in volumesList"))
			Expect(cueOutput).To(ContainSubstring("for j, vj in volumesList if j < i && vi.name == vj.name"))
			Expect(cueOutput).To(ContainSubstring("_ignore: true"))
			Expect(cueOutput).To(ContainSubstring("if val._ignore == _|_"))
			Expect(cueOutput).NotTo(ContainSubstring("[for v in volumesList { v }]"))
		})

		It("should create a Job via kube.#Apply with correct spec and conditional namespace", func() {
			Expect(cueOutput).To(ContainSubstring("job: kube.#Apply & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Job"`))
			Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepName)-\(context.stepSessionID)`))
			Expect(cueOutput).To(ContainSubstring(`parameter.serviceAccountName == "kubevela-vela-core"`))
			Expect(cueOutput).To(ContainSubstring(`namespace: "vela-system"`))
			Expect(cueOutput).To(ContainSubstring(`parameter.serviceAccountName != "kubevela-vela-core"`))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring("backoffLimit: 3"))
			Expect(cueOutput).To(ContainSubstring(`"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"`))
			Expect(cueOutput).To(ContainSubstring("image: parameter.image"))
			Expect(cueOutput).To(ContainSubstring("command: parameter.command"))
			Expect(cueOutput).To(ContainSubstring("volumeMounts: mountsArray"))
			Expect(cueOutput).To(ContainSubstring(`restartPolicy: "Never"`))
			Expect(cueOutput).To(ContainSubstring("serviceAccount: parameter.serviceAccountName"))
			Expect(cueOutput).To(ContainSubstring("volumes: deDupVolumesArray"))
		})

		It("should have log, fail, and wait actions for job lifecycle", func() {
			Expect(cueOutput).To(ContainSubstring("log: util.#Log & {"))
			Expect(cueOutput).To(ContainSubstring("labelSelector:"))
			Expect(cueOutput).To(ContainSubstring("job.$returns.value.status != _|_"))
			Expect(cueOutput).To(ContainSubstring("job.$returns.value.status.failed != _|_"))
			Expect(cueOutput).To(ContainSubstring("job.$returns.value.status.failed > 2"))
			Expect(cueOutput).To(ContainSubstring("breakWorkflow: builtin.#Fail & {"))
			Expect(cueOutput).To(ContainSubstring(`$params: message: "failed to execute vela command"`))
			Expect(cueOutput).To(ContainSubstring("wait: builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("$params: continue: job.$returns.value.status.succeeded > 0"))
		})

		It("should be structurally correct with expected action and loop counts", func() {
			Expect(strings.Count(cueOutput, "kube.#Apply & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "util.#Log & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#Fail & {")).To(Equal(1))

			mountsIdx := strings.Index(cueOutput, "mountsArray:")
			volumesIdx := strings.Index(cueOutput, "volumesList:")
			mountsSection := cueOutput[mountsIdx:volumesIdx]
			Expect(strings.Count(mountsSection, "for m in")).To(Equal(2))

			dedupIdx := strings.Index(cueOutput, "deDupVolumesArray:")
			volumesSection := cueOutput[volumesIdx:dedupIdx]
			Expect(strings.Count(volumesSection, "for m in")).To(Equal(2))
		})
	})
})
