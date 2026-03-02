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

var _ = Describe("Affinity Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Affinity()

		Expect(trait.GetName()).To(Equal("affinity"))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"ui-hidden": "true"`))

		// Parameters
		Expect(cue).To(ContainSubstring(`podAffinity?:`))
		Expect(cue).To(ContainSubstring(`podAntiAffinity?:`))
		Expect(cue).To(ContainSubstring(`nodeAffinity?:`))
		Expect(cue).To(ContainSubstring(`tolerations?:`))

		Expect(cue).To(ContainSubstring(`weight: int & >=1 & <=100`))

		Expect(cue).To(ContainSubstring(`podAffinityTerm: #podAffinityTerm`))
		Expect(cue).To(ContainSubstring(`nodeSelectorTerms: [...#nodeSelectorTerm]`))
		Expect(cue).To(ContainSubstring(`preference: #nodeSelectorTerm`))

		// Sub-field conditions
		Expect(cue).To(ContainSubstring(`parameter.podAffinity.required != _|_`))
		Expect(cue).To(ContainSubstring(`parameter.podAffinity.preferred != _|_`))
		Expect(cue).To(ContainSubstring(`parameter.podAntiAffinity.required != _|_`))
		Expect(cue).To(ContainSubstring(`parameter.nodeAffinity.required != _|_`))
		Expect(cue).To(ContainSubstring(`parameter.nodeAffinity.preferred != _|_`))

		// Optional field guards in foreach
		Expect(cue).To(ContainSubstring(`if v.labelSelector != _|_`))
		Expect(cue).To(ContainSubstring(`if v.namespaces != _|_`))
		Expect(cue).To(ContainSubstring(`if v.key != _|_`))
		Expect(cue).To(ContainSubstring(`if v.effect != _|_`))
		Expect(cue).To(ContainSubstring(`if v.tolerationSeconds != _|_`))
		Expect(cue).To(ContainSubstring(`operator: v.operator`)) // required field - no guard

		// Optional field guards for nested struct fields
		Expect(cue).To(ContainSubstring(`if v.namespaceSelector != _|_`))
		Expect(cue).To(ContainSubstring(`if v.matchExpressions != _|_`))
		Expect(cue).To(ContainSubstring(`if v.matchFields != _|_`))

		Expect(cue).To(ContainSubstring(`#labelSelector`))
		Expect(cue).To(ContainSubstring(`matchLabels?: [string]: string`))
		Expect(cue).To(ContainSubstring(`values?: [...string]`))
		Expect(cue).To(ContainSubstring(`namespaces?: [...string]`))
		Expect(cue).To(ContainSubstring(`#podAffinityTerm`))
		Expect(cue).To(ContainSubstring(`#nodeSelectorTerm`))
		Expect(cue).To(ContainSubstring(`matchExpressions?: [...#nodeSelector]`))
	})
})
