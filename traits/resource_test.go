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

var _ = Describe("Resource", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Resource()

		Expect(trait.GetName()).To(Equal("resource"))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"cronjobs.batch"`))

		// Parameters
		Expect(cue).To(ContainSubstring(`cpu?:`))
		Expect(cue).To(ContainSubstring(`memory?:`))
		Expect(cue).To(ContainSubstring(`*"2048Mi"`))
		Expect(cue).To(ContainSubstring(`=~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"`))
		Expect(cue).To(ContainSubstring(`requests?:`))
		Expect(cue).To(ContainSubstring(`limits?:`))

		// Template: let binding for DRY container element
		Expect(cue).To(ContainSubstring(`let resourceContent =`))
		Expect(cue).To(ContainSubstring(`containers: [resourceContent]`))

		// PatchStrategy annotations on requests/limits
		Expect(cue).To(ContainSubstring(`// +patchStrategy=retainKeys`))

		// Two-level context guards
		Expect(cue).To(ContainSubstring(`context.output.spec != _|_`))
		Expect(cue).To(ContainSubstring(`context.output.spec.template != _|_`))
		Expect(cue).To(ContainSubstring(`context.output.spec.jobTemplate != _|_`))
	})
})
