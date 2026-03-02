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

package traits_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/traits"
)

var _ = Describe("Env Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Env()

		Expect(trait.GetName()).To(Equal("env"))
		Expect(trait.GetDescription()).To(Equal("Add env on K8s pod for your workload which follows the pod spec in path 'spec.template'"))

		cue := trait.ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// #PatchParams schema: all 4 fields with correct types
		Expect(cue).To(ContainSubstring(`#PatchParams: {`))
		Expect(cue).To(ContainSubstring(`containerName: *"" | string`))
		Expect(cue).To(ContainSubstring(`replace: *false | bool`))
		Expect(cue).To(ContainSubstring(`env: [string]: string`))
		Expect(cue).To(ContainSubstring(`unset: *[] | [...string]`))

		// No duplicate containerName (1 in #PatchParams + 2 in _params mapping = 3)
		Expect(strings.Count(cue, "containerName:")).To(Equal(3))

		// PatchContainer body: complex env merge logic keywords
		Expect(cue).To(ContainSubstring(`PatchContainer: {`))
		Expect(cue).To(ContainSubstring(`_params: #PatchParams`))
		Expect(cue).To(ContainSubstring(`_delKeys: {for k in _params.unset`))
		Expect(cue).To(ContainSubstring(`_baseContainers: context.output.spec.template.spec.containers`))
		Expect(cue).To(ContainSubstring(`_baseEnv:       _baseContainer.env`))
		Expect(cue).To(ContainSubstring(`_baseEnvMap: {for envVar in _baseEnv`))
		Expect(cue).To(ContainSubstring(`envVar.valueFrom`))

		// _params mapping: auto-generated unconditional field mappings
		Expect(cue).To(ContainSubstring("replace: parameter.replace"))
		Expect(cue).To(ContainSubstring("env:     parameter.env"))
		Expect(cue).To(ContainSubstring("unset:   parameter.unset"))

		// Multi-container support
		Expect(cue).To(ContainSubstring("if parameter.containers == _|_"))
		Expect(cue).To(ContainSubstring("if parameter.containers != _|_"))
		Expect(cue).To(ContainSubstring("containers: [...#PatchParams]"))

		// Error collection
		Expect(cue).To(ContainSubstring(`errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`))

		// Descriptions
		Expect(cue).To(ContainSubstring("// +usage=Specify if replacing the whole environment settings"))
		Expect(cue).To(ContainSubstring("// +usage=Specify the  environment variables to merge"))
		Expect(cue).To(ContainSubstring("// +usage=Specify which existing environment variables to unset"))
		Expect(cue).To(ContainSubstring("// +usage=Specify the environment variables for multiple containers"))
	})
})
