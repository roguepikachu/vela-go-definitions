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

var _ = Describe("ServiceBinding Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.ServiceBinding()

		Expect(trait.GetName()).To(Equal("service-binding"))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`"ui-hidden": "true"`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))

		// Template: patch with patchKey annotations
		Expect(cue).To(ContainSubstring(`// +patchKey=name`))
		Expect(cue).To(ContainSubstring(`name: context.name`))

		// List comprehension over envMappings
		Expect(cue).To(ContainSubstring(`for envName, v in parameter.envMappings`))
		Expect(cue).To(ContainSubstring(`valueFrom: secretKeyRef:`))
		Expect(cue).To(ContainSubstring(`if v["key"] != _|_`))
		Expect(cue).To(ContainSubstring(`if v["key"] == _|_`))

		// Fluent parameter
		Expect(cue).To(ContainSubstring(`envMappings: [string]: #KeySecret`))

		// Fluent helper definition
		Expect(cue).To(ContainSubstring(`#KeySecret:`))
		Expect(cue).To(ContainSubstring(`key?:`))
		Expect(cue).To(ContainSubstring(`secret: string`))
	})
})
