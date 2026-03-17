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

var _ = Describe("K8sUpdateStrategy", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.K8sUpdateStrategy()

		Expect(trait.GetName()).To(Equal("k8s-update-strategy"))
		Expect(trait.GetDescription()).To(Equal("Set k8s update strategy for Deployment/DaemonSet/StatefulSet"))

		cue := trait.ToCue()

		// Three separate conditional blocks for each workload type
		Expect(cue).To(ContainSubstring(`parameter.targetKind == "Deployment" && parameter.strategy.type != "OnDelete"`))
		Expect(cue).To(ContainSubstring(`parameter.targetKind == "StatefulSet" && parameter.strategy.type != "Recreate"`))
		Expect(cue).To(ContainSubstring(`parameter.targetKind == "DaemonSet" && parameter.strategy.type != "Recreate"`))

		// Three patchStrategy annotations
		Expect(strings.Count(cue, "// +patchStrategy=retainKeys")).To(Equal(3))

		// Deployment uses "strategy", StatefulSet/DaemonSet use "updateStrategy"
		Expect(cue).To(ContainSubstring("strategy: {"))
		Expect(cue).To(ContainSubstring("updateStrategy: {"))

		// Inner RollingUpdate condition
		Expect(cue).To(ContainSubstring(`parameter.strategy.type == "RollingUpdate"`))

		// Correct field assignments
		Expect(cue).To(ContainSubstring("maxSurge:       parameter.strategy.rollingStrategy.maxSurge"))
		Expect(cue).To(ContainSubstring("maxUnavailable: parameter.strategy.rollingStrategy.maxUnavailable"))
		Expect(cue).To(ContainSubstring("partition: parameter.strategy.rollingStrategy.partition"))

		// Parameters
		Expect(cue).To(ContainSubstring(`targetAPIVersion: *"apps/v1" | string`))
		Expect(cue).To(ContainSubstring(`targetKind: *"Deployment" | "StatefulSet" | "DaemonSet"`))
		Expect(cue).To(ContainSubstring(`type: *"RollingUpdate" | "Recreate" | "OnDelete"`))
	})

	It("should have optional rollingStrategy field", func() {
		trait := traits.K8sUpdateStrategy()
		cue := trait.ToCue()

		// rollingStrategy must be optional (not needed for Recreate strategy)
		Expect(cue).To(ContainSubstring(`rollingStrategy?: {`))
	})
})
