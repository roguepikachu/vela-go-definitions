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

var _ = Describe("StartupProbe Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.StartupProbe()

		Expect(trait.GetName()).To(Equal("startup-probe"))
		Expect(trait.GetDescription()).To(Equal("Add startup probe hooks for the specified container of K8s pod for your workload which follows the pod spec in path 'spec.template'."))

		cue := trait.ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// #PatchParams: fields with defaults use *default | type
		Expect(cue).To(ContainSubstring(`initialDelaySeconds: *0 | int`))
		Expect(cue).To(ContainSubstring(`periodSeconds: *10 | int`))
		Expect(cue).To(ContainSubstring(`timeoutSeconds: *1 | int`))
		Expect(cue).To(ContainSubstring(`successThreshold: *1 | int`))
		Expect(cue).To(ContainSubstring(`failureThreshold: *3 | int`))

		// #PatchParams: optional fields use field?: type syntax
		Expect(cue).To(ContainSubstring(`terminationGracePeriodSeconds?: int`))
		Expect(cue).To(ContainSubstring(`exec?: {`))
		Expect(cue).To(ContainSubstring(`httpGet?: {`))
		Expect(cue).To(ContainSubstring(`grpc?: {`))
		Expect(cue).To(ContainSubstring(`tcpSocket?: {`))

		// PatchContainer structure
		Expect(cue).To(ContainSubstring(`#StartupProbeParams: {`))
		Expect(cue).To(ContainSubstring(`PatchContainer: {`))
		Expect(cue).To(ContainSubstring(`_params:         #StartupProbeParams`))
		Expect(cue).To(ContainSubstring(`_baseContainers: context.output.spec.template.spec.containers`))

		// PatchContainer body: conditional blocks for optional probe types
		Expect(cue).To(ContainSubstring(`if _params.exec != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.httpGet != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.grpc != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.tcpSocket != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.terminationGracePeriodSeconds != _|_`))

		// PatchContainer body: conditional blocks for fields with IsSet().Default()
		Expect(cue).To(ContainSubstring(`if _params.initialDelaySeconds != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.periodSeconds != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.timeoutSeconds != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.successThreshold != _|_`))
		Expect(cue).To(ContainSubstring(`if _params.failureThreshold != _|_`))

		// startupProbe group wrapper
		Expect(cue).To(ContainSubstring(`startupProbe: {`))

		// Multi-container support with custom param name "probes"
		Expect(cue).To(ContainSubstring(`parameter: *#StartupProbeParams | close({`))
		Expect(cue).To(ContainSubstring(`probes: [...#StartupProbeParams]`))
		Expect(cue).To(ContainSubstring(`// +usage=Specify the startup probe for multiple containers`))

		// Error collection
		Expect(cue).To(ContainSubstring(`errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`))

		// Descriptions for probe fields
		Expect(cue).To(ContainSubstring(`// +usage=Number of seconds after the container has started before liveness probes are initiated`))
		Expect(cue).To(ContainSubstring(`// +usage=How often, in seconds, to execute the probe`))
		Expect(cue).To(ContainSubstring(`// +usage=Number of seconds after which the probe times out`))

		// No duplicate parameter blocks
		Expect(strings.Count(cue, "parameter:")).To(Equal(1))
	})
})
