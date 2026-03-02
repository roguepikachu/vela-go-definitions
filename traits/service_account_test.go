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

var _ = Describe("ServiceAccount Trait", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.ServiceAccount()

		Expect(trait.GetName()).To(Equal("service-account"))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: false`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Let bindings for filtered privilege arrays
		Expect(cue).To(ContainSubstring(`let _clusterPrivileges =`))
		Expect(cue).To(ContainSubstring(`let _namespacePrivileges =`))
		Expect(cue).To(ContainSubstring(`v.scope == "cluster"`))
		Expect(cue).To(ContainSubstring(`v.scope == "namespace"`))

		// Patch
		Expect(cue).To(ContainSubstring(`// +patchStrategy=retainKeys`))
		Expect(cue).To(ContainSubstring(`serviceAccountName: parameter.name`))

		// Conditional ServiceAccount output
		Expect(cue).To(ContainSubstring(`if parameter.create`))
		Expect(cue).To(ContainSubstring(`"service-account":`))
		Expect(cue).To(ContainSubstring(`kind:       "ServiceAccount"`))

		// Conditional cluster-scoped RBAC output group
		Expect(cue).To(ContainSubstring(`len(_clusterPrivileges) > 0`))
		Expect(cue).To(ContainSubstring(`"cluster-role":`))
		Expect(cue).To(ContainSubstring(`kind:       "ClusterRole"`))
		Expect(cue).To(ContainSubstring(`"cluster-role-binding":`))
		Expect(cue).To(ContainSubstring(`kind:       "ClusterRoleBinding"`))

		// Conditional namespace-scoped RBAC output group
		Expect(cue).To(ContainSubstring(`len(_namespacePrivileges) > 0`))
		Expect(cue).To(ContainSubstring(`kind:       "Role"`))
		Expect(cue).To(ContainSubstring(`kind:       "RoleBinding"`))

		// String interpolation for cluster-scoped resource names
		Expect(cue).To(ContainSubstring(`"\(context.namespace):\(parameter.name)"`))

		// Rules comprehension with optional fields
		Expect(cue).To(ContainSubstring(`for v in _clusterPrivileges`))
		Expect(cue).To(ContainSubstring(`verbs: v.verbs`))
		Expect(cue).To(ContainSubstring(`if v.apiGroups != _|_`))

		// Helper type definition
		Expect(cue).To(ContainSubstring(`#Privileges`))
		Expect(cue).To(ContainSubstring(`privileges?: [...#Privileges]`))
		Expect(cue).To(ContainSubstring(`scope: *"namespace" | "cluster"`))

		// Parameters
		Expect(cue).To(ContainSubstring(`name: string`))
		Expect(cue).To(ContainSubstring(`create: *false | bool`))
	})
})
