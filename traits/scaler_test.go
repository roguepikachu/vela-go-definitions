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

var _ = Describe("Scaler Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Scaler()

		Expect(trait.GetName()).To(Equal("scaler"))
		Expect(trait.GetDescription()).To(Equal("Manually scale K8s pod for your workload which follows the pod spec in path 'spec.template'."))

		cue := trait.ToCue()

		// Verify trait metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))

		// Verify replicas parameter has correct type and default
		Expect(cue).To(ContainSubstring(`replicas: *1 | int`))

		// Verify patch targets spec.replicas with retainKeys strategy
		Expect(cue).To(ContainSubstring(`// +patchStrategy=retainKeys`))
		Expect(cue).To(ContainSubstring(`spec: replicas: parameter.replicas`))
	})
})
