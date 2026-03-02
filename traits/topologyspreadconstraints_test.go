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

var _ = Describe("TopologySpreadConstraints Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.TopologySpreadConstraints()

		Expect(trait.GetName()).To(Equal("topologyspreadconstraints"))

		cue := trait.ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Bug 1 fix: constraints array should be required (no ?)
		Expect(cue).To(ContainSubstring("constraints: [...{"))
		Expect(cue).NotTo(ContainSubstring("constraints?: [...{"))

		// Bug 2 fix: labelSelector should be required and reference #labSelector helper
		Expect(cue).To(ContainSubstring("labelSelector: #labSelector"))

		// Helper type definition should exist as closed struct
		Expect(cue).To(ContainSubstring("#labSelector: {"))

		// Bug 3 fix: nodeAffinityPolicy and nodeTaintsPolicy should be optional WITH default
		Expect(cue).To(ContainSubstring(`nodeAffinityPolicy?: *"Honor" | "Ignore"`))
		Expect(cue).To(ContainSubstring(`nodeTaintsPolicy?: *"Honor" | "Ignore"`))
		// Must NOT be required (without ?)
		Expect(cue).NotTo(ContainSubstring(`nodeAffinityPolicy: *"Honor"`))
		Expect(cue).NotTo(ContainSubstring(`nodeTaintsPolicy: *"Honor"`))

		// Other parameter fields
		Expect(cue).To(ContainSubstring(`maxSkew: int`))
		Expect(cue).To(ContainSubstring(`topologyKey: string`))
		Expect(cue).To(ContainSubstring(`whenUnsatisfiable: *"DoNotSchedule" | "ScheduleAnyway"`))
		Expect(cue).To(ContainSubstring(`minDomains?: int`))
		Expect(cue).To(ContainSubstring(`matchLabelKeys?: [...string]`))
		Expect(cue).To(ContainSubstring(`matchLabels?: [string]: string`))
		Expect(cue).To(ContainSubstring(`operator: *"In" | "NotIn" | "Exists" | "DoesNotExist"`))
		Expect(cue).To(ContainSubstring(`values?: [...string]`))

		// Template: conditional field guards for optional fields
		Expect(cue).To(ContainSubstring(`if v.nodeAffinityPolicy != _|_`))
		Expect(cue).To(ContainSubstring(`if v.nodeTaintsPolicy != _|_`))
		Expect(cue).To(ContainSubstring(`if v.minDomains != _|_`))
		Expect(cue).To(ContainSubstring(`if v.matchLabelKeys != _|_`))
	})
})
