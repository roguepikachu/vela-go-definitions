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

var _ = Describe("Gateway Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Gateway()

		Expect(trait.GetName()).To(Equal("gateway"))
		Expect(trait.GetDescription()).To(Equal("Enable public web traffic for the component, the ingress API matches K8s v1.20+."))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`customStatus:`))
		Expect(cue).To(ContainSubstring(`healthPolicy:`))

		// Import
		Expect(cue).To(ContainSubstring(`"strconv"`))

		// Let bindings with conditional values
		Expect(cue).To(ContainSubstring(`let nameSuffix =`))
		Expect(cue).To(ContainSubstring(`let serviceMetaName =`))

		// Conditional Service output (only when no existing service)
		Expect(cue).To(ContainSubstring(`if (parameter.existingServiceName == _|_)`))
		Expect(cue).To(ContainSubstring(`kind:       "Service"`))

		// Dynamic output names
		Expect(cue).To(ContainSubstring(`(serviceOutputName):`))
		Expect(cue).To(ContainSubstring(`(ingressOutputName):`))

		// Cluster version conditional apiVersion for Ingress
		Expect(cue).To(ContainSubstring(`legacyAPI:`))
		Expect(cue).To(ContainSubstring(`context.clusterVersion.minor < 19`))
		Expect(cue).To(ContainSubstring(`"networking.k8s.io/v1beta1"`))
		Expect(cue).To(ContainSubstring(`"networking.k8s.io/v1"`))
		Expect(cue).To(ContainSubstring(`kind: "Ingress"`))

		// Map iteration for ports and paths
		Expect(cue).To(ContainSubstring(`for k, v in parameter.http`))
		Expect(cue).To(ContainSubstring(`strconv.FormatInt`))

		// Conditional annotations and labels spreading
		Expect(cue).To(ContainSubstring(`if parameter.annotations != _|_`))
		Expect(cue).To(ContainSubstring(`for key, value in parameter.annotations`))
		Expect(cue).To(ContainSubstring(`if parameter.labels != _|_`))

		// Parameters
		Expect(cue).To(ContainSubstring(`domain?: string`))
		Expect(cue).To(ContainSubstring(`http: [string]: int`))
		Expect(cue).To(ContainSubstring(`class: *"nginx" | string`))
		Expect(cue).To(ContainSubstring(`classInSpec: *false | bool`))
		Expect(cue).To(ContainSubstring(`secretName?: string`))
		Expect(cue).To(ContainSubstring(`pathType: *"ImplementationSpecific"`))
		Expect(cue).To(ContainSubstring(`existingServiceName?: string`))
	})
})
