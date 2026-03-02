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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/traits"
)

var _ = Describe("SecurityContext Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.SecurityContext()

		Expect(trait.GetName()).To(Equal("securitycontext"))
		Expect(trait.GetDescription()).To(Equal("Adds security context to the container spec in path 'spec.template.spec.containers.[].securityContext'."))

		cue := trait.ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// #PatchParams: fields with explicit defaults use *default | type
		Expect(cue).To(ContainSubstring(`containerName: *"" | string`))
		Expect(cue).To(ContainSubstring(`allowPrivilegeEscalation: *false | bool`))
		Expect(cue).To(ContainSubstring(`readOnlyRootFilesystem: *false | bool`))
		Expect(cue).To(ContainSubstring(`privileged: *false | bool`))
		Expect(cue).To(ContainSubstring(`runAsNonRoot: *true | bool`))

		// #PatchParams: fields with != _|_ condition use optional syntax (field?: type)
		Expect(cue).To(ContainSubstring(`runAsUser?: int`))
		Expect(cue).To(ContainSubstring(`runAsGroup?: int`))
		Expect(cue).To(ContainSubstring(`addCapabilities?: [...string]`))
		Expect(cue).To(ContainSubstring(`dropCapabilities?: [...string]`))

		// Must NOT have *null | type for optional fields
		Expect(cue).NotTo(ContainSubstring(`runAsUser: *null | int`))
		Expect(cue).NotTo(ContainSubstring(`runAsGroup: *null | int`))
		Expect(cue).NotTo(ContainSubstring(`addCapabilities: *null`))
		Expect(cue).NotTo(ContainSubstring(`dropCapabilities: *null`))

		// PatchContainer structure
		Expect(cue).To(ContainSubstring(`#PatchParams: {`))
		Expect(cue).To(ContainSubstring(`PatchContainer: {`))
		Expect(cue).To(ContainSubstring(`_params:         #PatchParams`))

		// PatchContainer body: conditional blocks for optional fields
		Expect(cue).To(ContainSubstring(`if _params.runAsUser != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.runAsGroup != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.addCapabilities != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.dropCapabilities != _|_`))

		// PatchContainer body: unconditional assignments for fields with defaults
		Expect(cue).To(ContainSubstring(`allowPrivilegeEscalation: _params.allowPrivilegeEscalation`))
		Expect(cue).To(ContainSubstring(`readOnlyRootFilesystem:   _params.readOnlyRootFilesystem`))
		Expect(cue).To(ContainSubstring(`privileged:               _params.privileged`))
		Expect(cue).To(ContainSubstring(`runAsNonRoot:             _params.runAsNonRoot`))

		// Multi-container support
		Expect(cue).To(ContainSubstring("parameter: #PatchParams | close({"))
		Expect(cue).To(ContainSubstring("containers: [...#PatchParams]"))

		// Error collection
		Expect(cue).To(ContainSubstring(`errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`))
	})
})
