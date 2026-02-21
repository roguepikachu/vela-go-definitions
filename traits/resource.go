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

package traits

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Resource creates the resource trait definition.
// This trait adds resource requests and limits on K8s pods.
//
// The template uses SetRawPatchBlock because it requires:
// - A let binding (let resourceContent) to DRY the container element across two patch paths
// - patchStrategy=retainKeys annotations on requests/limits fields
// - Chained if conditions for two-level context guards (spec != _|_ if spec.template != _|_)
// Parameters are still defined fluently.
func Resource() *defkit.TraitDefinition {
	// Shorthand parameters for simple cases - using custom schema for union types
	cpu := defkit.Map("cpu").WithSchema(`*1 | number | string`).Description("Specify the amount of cpu for requests and limits")
	memory := defkit.Map("memory").WithSchema(`*"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"`).Description("Specify the amount of memory for requests and limits")

	// Explicit requests parameter - using custom schema for the structured type
	requests := defkit.Map("requests").Description("Specify the resources in requests").WithSchema(`{
		// +usage=Specify the amount of cpu for requests
		cpu: *1 | number | string
		// +usage=Specify the amount of memory for requests
		memory: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
	}`)

	// Explicit limits parameter
	limits := defkit.Map("limits").Description("Specify the resources in limits").WithSchema(`{
		// +usage=Specify the amount of cpu for limits
		cpu: *1 | number | string
		// +usage=Specify the amount of memory for limits
		memory: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
	}`)

	return defkit.NewTrait("resource").
		Description("Add resource requests and limits on K8s pod for your workload which follows the pod spec in path 'spec.template.'").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "cronjobs.batch").
		PodDisruptive(true).
		Params(cpu, memory, requests, limits).
		Template(func(tpl *defkit.Template) {
			tpl.SetRawPatchBlock(`patch: {
	let resourceContent = {
		resources: {
			if parameter.cpu != _|_ if parameter.memory != _|_ if parameter.requests == _|_ if parameter.limits == _|_ {
				// +patchStrategy=retainKeys
				requests: {
					cpu:    parameter.cpu
					memory: parameter.memory
				}
				// +patchStrategy=retainKeys
				limits: {
					cpu:    parameter.cpu
					memory: parameter.memory
				}
			}

			if parameter.requests != _|_ {
				// +patchStrategy=retainKeys
				requests: {
					cpu:    parameter.requests.cpu
					memory: parameter.requests.memory
				}
			}
			if parameter.limits != _|_ {
				// +patchStrategy=retainKeys
				limits: {
					cpu:    parameter.limits.cpu
					memory: parameter.limits.memory
				}
			}
		}
	}

	if context.output.spec != _|_ if context.output.spec.template != _|_ {
		spec: template: spec: {
			// +patchKey=name
			containers: [resourceContent]
		}
	}
	if context.output.spec != _|_ if context.output.spec.jobTemplate != _|_ {
		spec: jobTemplate: spec: template: spec: {
			// +patchKey=name
			containers: [resourceContent]
		}
	}
}`)
		})
}

func init() {
	defkit.Register(Resource())
}
