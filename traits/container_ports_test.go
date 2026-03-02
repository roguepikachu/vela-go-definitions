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

var _ = Describe("Container Ports Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.ContainerPorts()

		Expect(trait.GetName()).To(Equal("container-ports"))
		Expect(trait.GetDescription()).To(Equal("Expose on the host and bind the external port to host to enable web traffic for your component."))

		cue := trait.ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Imports
		Expect(cue).To(ContainSubstring(`"strconv"`))
		Expect(cue).To(ContainSubstring(`"strings"`))

		// #PatchParams schema: containerName + ports with nested struct
		Expect(cue).To(ContainSubstring(`#PatchParams: {`))
		Expect(cue).To(ContainSubstring(`containerName: *"" | string`))
		Expect(cue).To(ContainSubstring(`ports: *[] | [...{`))
		Expect(cue).To(ContainSubstring(`containerPort: int`))
		Expect(cue).To(ContainSubstring(`protocol: *"TCP" | "UDP" | "SCTP"`))
		Expect(cue).To(ContainSubstring(`hostPort?: int`))
		Expect(cue).To(ContainSubstring(`hostIP?: string`))

		// No duplicate containerName (1 in #PatchParams + 2 in _params mapping = 3)
		Expect(strings.Count(cue, "containerName:")).To(Equal(3))

		// PatchContainer body: complex port merge logic
		Expect(cue).To(ContainSubstring(`PatchContainer: {`))
		Expect(cue).To(ContainSubstring(`_params:         #PatchParams`))
		Expect(cue).To(ContainSubstring(`_baseContainers: context.output.spec.template.spec.containers`))
		Expect(cue).To(ContainSubstring(`_basePorts:     _baseContainer.ports`))
		Expect(cue).To(ContainSubstring(`_basePortsMap:`))
		Expect(cue).To(ContainSubstring(`_portsMap:`))
		Expect(cue).To(ContainSubstring(`_uniqueKey:`))
		Expect(cue).To(ContainSubstring(`strings.ToLower`))
		Expect(cue).To(ContainSubstring(`strconv.FormatInt`))

		// _params mapping: auto-generated
		Expect(cue).To(ContainSubstring("ports: parameter.ports"))

		// Multi-container support
		Expect(cue).To(ContainSubstring("if parameter.containers == _|_"))
		Expect(cue).To(ContainSubstring("if parameter.containers != _|_"))
		Expect(cue).To(ContainSubstring("containers: [...#PatchParams]"))

		// Error collection
		Expect(cue).To(ContainSubstring(`errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`))

		// Descriptions
		Expect(cue).To(ContainSubstring("// +usage=Specify ports you want customer traffic sent to"))
		Expect(cue).To(ContainSubstring("// +usage=Specify the container ports for multiple containers"))
	})
})
