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

var _ = Describe("Annotations Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Annotations()

		Expect(trait.GetName()).To(Equal("annotations"))

		cue := trait.ToCue()

		// Verify raw CUE content is present
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`metadata: annotations:`))
		Expect(cue).To(ContainSubstring(`context.output.spec`))
		Expect(cue).To(ContainSubstring(`jobTemplate`))
		Expect(cue).To(ContainSubstring(`parameter: [string]: string | null`))

		// Let binding: annotationsContent should be defined once with ForEachMap
		Expect(cue).To(ContainSubstring(`let annotationsContent =`))
		Expect(cue).To(ContainSubstring(`for k, v in parameter`))
		Expect(cue).To(ContainSubstring(`(k): v`))

		// The for-each comprehension should appear only once (in the let binding),
		// not inlined at each of the 4 usage sites
		Expect(strings.Count(cue, "for k, v in parameter")).To(Equal(1))

		// All 4 annotation sites should reference the let variable
		Expect(cue).To(ContainSubstring(`metadata: annotations: annotationsContent`))
		Expect(cue).To(ContainSubstring(`annotations: annotationsContent`))

		// Conditional blocks for spec.template and jobTemplate should reference let variable
		Expect(cue).To(ContainSubstring(`context.output.spec.template != _|_`))
		Expect(cue).To(ContainSubstring(`context.output.spec.jobTemplate != _|_`))
		Expect(cue).To(ContainSubstring(`context.output.spec.jobTemplate.spec != _|_`))
		Expect(cue).To(ContainSubstring(`context.output.spec.jobTemplate.spec.template != _|_`))

		// Count references to annotationsContent (1 let definition + 4 usage sites = 5)
		Expect(strings.Count(cue, "annotationsContent")).To(Equal(5))
	})
})
