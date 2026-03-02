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

var _ = Describe("Sidecar Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Sidecar()

		Expect(trait.GetName()).To(Equal("sidecar"))
		Expect(trait.GetDescription()).To(Equal("Inject a sidecar container to K8s pod for your workload which follows the pod spec in path 'spec.template'."))

		cue := trait.ToCue()

		// Verify trait metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Verify required parameters with types
		Expect(cue).To(ContainSubstring(`name: string`))
		Expect(cue).To(ContainSubstring(`image: string`))

		// Verify optional sidecar parameters
		Expect(cue).To(ContainSubstring(`cmd?: [...string]`))
		Expect(cue).To(ContainSubstring(`args?: [...string]`))
		Expect(cue).To(ContainSubstring(`env?: [...{`))
		Expect(cue).To(ContainSubstring(`volumes?: [...{`))

		// Verify sidecar container injection via patchKey
		Expect(cue).To(ContainSubstring(`// +patchKey=name`))
		Expect(cue).To(ContainSubstring(`containers:`))

		// Verify HealthProbe reference for probes
		Expect(cue).To(ContainSubstring(`#HealthProbe`))
		Expect(cue).To(ContainSubstring(`livenessProbe?:`))
		Expect(cue).To(ContainSubstring(`readinessProbe?:`))

		// #HealthProbe exec.command should have string element type
		Expect(cue).To(ContainSubstring(`command: [...string]`))
		Expect(cue).NotTo(ContainSubstring("command: [...]"))

		// #HealthProbe httpGet.httpHeaders should have structured elements
		Expect(cue).To(ContainSubstring(`httpHeaders?: [...{`))
		Expect(cue).To(ContainSubstring("name:  string"))
		Expect(cue).To(ContainSubstring("value: string"))
	})
})
