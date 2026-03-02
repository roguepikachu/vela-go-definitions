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

var _ = Describe("HostAlias", func() {
	It("should have correct name and CUE output", func() {
		cue := traits.HostAlias().ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`hostalias: {`))
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`description: "Add host aliases on K8s pod for your workload which follows the pod spec in path 'spec.template'."`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Patch block: patchKey annotation and direct array assignment (no wrapping)
		Expect(cue).To(ContainSubstring(`// +patchKey=ip`))
		Expect(cue).To(ContainSubstring(`hostAliases: parameter.hostAliases`))
		// Should NOT wrap in array brackets
		Expect(cue).NotTo(ContainSubstring(`[parameter.hostAliases]`))

		// Parameter block: hostAliases should be required (no ?)
		Expect(cue).To(ContainSubstring("hostAliases: [...{"))
		Expect(cue).NotTo(ContainSubstring("hostAliases?: [...{"))

		// Struct fields inside hostAliases
		Expect(cue).To(ContainSubstring(`ip: string`))
		Expect(cue).To(ContainSubstring(`hostnames: [...string]`))

		// Description
		Expect(cue).To(ContainSubstring(`// +usage=Specify the hostAliases to add`))
	})
})
