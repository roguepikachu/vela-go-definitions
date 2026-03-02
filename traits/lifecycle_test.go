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

var _ = Describe("Lifecycle", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Lifecycle()

		Expect(trait.GetName()).To(Equal("lifecycle"))
		Expect(trait.GetDescription()).To(Equal("Add lifecycle hooks for every container of K8s pod for your workload which follows the pod spec in path 'spec.template'."))

		cue := trait.ToCue()

		// Metadata
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Spread constraint pattern [...{struct}] for containers
		Expect(cue).To(ContainSubstring(`containers: [...{`))
		Expect(cue).To(ContainSubstring(`lifecycle: {`))
		Expect(cue).To(ContainSubstring(`if parameter["postStart"] != _|_ {`))
		Expect(cue).To(ContainSubstring(`postStart: parameter.postStart`))
		Expect(cue).To(ContainSubstring(`if parameter["preStop"] != _|_ {`))
		Expect(cue).To(ContainSubstring(`preStop: parameter.preStop`))
		Expect(cue).NotTo(ContainSubstring(`+patchKey`))

		// #Port is a constrained int
		Expect(cue).To(ContainSubstring(`#Port: int & >=1 & <=65535`))

		// Port fields reference #Port helper
		Expect(cue).To(ContainSubstring(`port:   #Port`))
		lines := strings.Split(cue, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "port") && strings.Contains(trimmed, "int") && !strings.Contains(trimmed, "#Port") {
				Expect(trimmed).NotTo(And(ContainSubstring("int"), Not(ContainSubstring("#Port"))),
					"Found port field without #Port reference: "+trimmed)
			}
		}

		// exec command is typed [...string] and required
		Expect(cue).To(ContainSubstring(`exec?: command: [...string]`))
		Expect(cue).NotTo(ContainSubstring(`command?: [...string]`))

		// httpHeaders is typed struct array
		Expect(cue).To(ContainSubstring(`httpHeaders?: [...{`))
		Expect(cue).To(ContainSubstring(`name:  string`))
		Expect(cue).To(ContainSubstring(`value: string`))

		// Parameters reference #LifeCycleHandler
		Expect(cue).To(ContainSubstring(`postStart?: #LifeCycleHandler`))
		Expect(cue).To(ContainSubstring(`preStop?: #LifeCycleHandler`))

		// Helper definitions
		Expect(cue).To(ContainSubstring(`#LifeCycleHandler: {`))
		Expect(cue).To(ContainSubstring(`scheme: *"HTTP" | "HTTPS"`))
		Expect(cue).To(ContainSubstring(`tcpSocket?: {`))
		Expect(cue).To(ContainSubstring(`host?: string`))
	})
})
