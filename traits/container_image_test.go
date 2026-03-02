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

var _ = Describe("Container Image Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.ContainerImage()

		Expect(trait.GetName()).To(Equal("container-image"))
		Expect(trait.GetDescription()).To(Equal("Set the image of the container."))

		cue := trait.ToCue()

		// Verify trait metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))

		Expect(cue).To(ContainSubstring(`imagePullPolicy: *"" | "IfNotPresent" | "Always" | "Never"`))
		Expect(cue).NotTo(ContainSubstring(`imagePullPolicy: *null`))

		Expect(cue).To(ContainSubstring("imagePullPolicy: parameter.imagePullPolicy"))

		Expect(cue).To(ContainSubstring("parameter: #PatchParams | close({"))

		Expect(strings.Count(cue, "parameter:")).To(Equal(1))
		Expect(cue).NotTo(ContainSubstring("parameter: {}"))

		Expect(cue).To(ContainSubstring("// +usage=Specify the image of the container"))
		Expect(cue).To(ContainSubstring("// +usage=Specify the image pull policy of the container"))
		Expect(cue).To(ContainSubstring("// +usage=Specify the container image for multiple containers"))

		// PatchContainer structure
		Expect(cue).To(ContainSubstring(`#PatchParams: {`))
		Expect(cue).To(ContainSubstring(`PatchContainer: {`))
		Expect(cue).To(ContainSubstring(`_params:         #PatchParams`))
		Expect(cue).To(ContainSubstring(`_baseContainers: context.output.spec.template.spec.containers`))
		Expect(cue).To(ContainSubstring(`errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`))

		// PatchContainer body: conditional for imagePullPolicy inside PatchContainer
		Expect(cue).To(ContainSubstring(`if _params.imagePullPolicy != ""`))

		// Multi-container support
		Expect(cue).To(ContainSubstring("if parameter.containers == _|_"))
		Expect(cue).To(ContainSubstring("if parameter.containers != _|_"))
		Expect(cue).To(ContainSubstring("containers: [...#PatchParams]"))
	})
})
