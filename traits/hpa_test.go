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

var _ = Describe("HPA", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.HPA()

		Expect(trait.GetName()).To(Equal("hpa"))
		Expect(trait.GetDescription()).To(Equal("Configure k8s HPA for Deployment or Statefulsets"))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))

		// Conditional apiVersion based on cluster version
		Expect(cue).To(ContainSubstring(`if context.clusterVersion.minor < 23`))
		Expect(cue).To(ContainSubstring(`apiVersion: "autoscaling/v2beta2"`))
		Expect(cue).To(ContainSubstring(`if context.clusterVersion.minor >= 23`))
		Expect(cue).To(ContainSubstring(`apiVersion: "autoscaling/v2"`))

		// Output resource
		Expect(cue).To(ContainSubstring(`outputs: hpa:`))
		Expect(cue).To(ContainSubstring(`kind: "HorizontalPodAutoscaler"`))
		Expect(cue).To(ContainSubstring(`metadata: name: context.name`))

		// Scale target ref
		Expect(cue).To(ContainSubstring(`scaleTargetRef:`))
		Expect(cue).To(ContainSubstring(`parameter.targetAPIVersion`))
		Expect(cue).To(ContainSubstring(`parameter.targetKind`))

		// Metrics array: static CPU, conditional memory, iterated custom
		Expect(cue).To(ContainSubstring(`metrics:`))
		Expect(cue).To(ContainSubstring(`name: "cpu"`))
		Expect(cue).To(ContainSubstring(`if parameter["mem"] != _|_`))
		Expect(cue).To(ContainSubstring(`name: "memory"`))
		Expect(cue).To(ContainSubstring(`if parameter["podCustomMetrics"] != _|_ for m in parameter.podCustomMetrics`))
		Expect(cue).To(ContainSubstring(`type: "Pods"`))

		// Conditional target type for CPU/memory
		Expect(cue).To(ContainSubstring(`if parameter.cpu.type == "Utilization"`))
		Expect(cue).To(ContainSubstring(`averageUtilization: parameter.cpu.value`))
		Expect(cue).To(ContainSubstring(`if parameter.cpu.type == "AverageValue"`))
		Expect(cue).To(ContainSubstring(`averageValue: parameter.cpu.value`))

		// Parameters
		Expect(cue).To(ContainSubstring(`min: *1 | int`))
		Expect(cue).To(ContainSubstring(`max: *10 | int`))
		Expect(cue).To(ContainSubstring(`targetAPIVersion: *"apps/v1" | string`))
		Expect(cue).To(ContainSubstring(`targetKind: *"Deployment" | string`))
		Expect(cue).To(ContainSubstring(`mem?:`))
		Expect(cue).To(ContainSubstring(`podCustomMetrics?:`))
	})
})
