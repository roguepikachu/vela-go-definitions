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

var _ = Describe("Expose Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Expose()

		Expect(trait.GetName()).To(Equal("expose"))
		Expect(trait.GetDescription()).To(Equal("Expose port to enable web traffic for your component."))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))
		Expect(cue).To(ContainSubstring(`stage:`))
		Expect(cue).To(ContainSubstring(`"PostDispatch"`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`customStatus:`))
		Expect(cue).To(ContainSubstring(`healthPolicy:`))

		// Imports
		Expect(cue).To(ContainSubstring(`"strconv"`))
		Expect(cue).To(ContainSubstring(`"strings"`))

		// Output resource
		Expect(cue).To(ContainSubstring(`outputs: service:`))
		Expect(cue).To(ContainSubstring(`kind:       "Service"`))
		Expect(cue).To(ContainSubstring(`metadata: name:        context.name`))

		// Dual-path port handling (legacy vs modern)
		Expect(cue).To(ContainSubstring(`if parameter["port"] != _|_`))
		Expect(cue).To(ContainSubstring(`if parameter["ports"] != _|_`))
		Expect(cue).To(ContainSubstring(`strconv.FormatInt`))
		Expect(cue).To(ContainSubstring(`strings.ToLower`))

		// Parameters
		Expect(cue).To(ContainSubstring(`port?: [...int]`))
		Expect(cue).To(ContainSubstring(`ports?: [`))
		Expect(cue).To(ContainSubstring(`annotations: [string]:`))
		Expect(cue).To(ContainSubstring(`matchLabels?: [string]:`))
		Expect(cue).To(ContainSubstring(`*"ClusterIP"`))
	})
})
