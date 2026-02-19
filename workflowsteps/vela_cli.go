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

package workflowsteps

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// VelaCli creates the vela-cli workflow step definition.
// This step runs a vela command inside a Kubernetes Job.
func VelaCli() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()
	stepName := defkit.Reference("context.stepName")
	stepSessionID := defkit.Reference("context.stepSessionID")

	// Parameters
	command := defkit.StringList("command").
		Description("Specify the vela command")
	image := defkit.String("image").
		Default("oamdev/vela-cli:v1.6.4").
		Description("Specify the image")
	serviceAccountName := defkit.String("serviceAccountName").
		Default("kubevela-vela-core").
		Description("specify serviceAccountName want to use")
	storage := defkit.Struct("storage").
		Fields(
			defkit.Field("secret", defkit.ParamTypeArray).
				ArrayOf(defkit.ParamTypeStruct).
				Nested(defkit.Struct("").Fields(
					defkit.Field("name", defkit.ParamTypeString).Required(),
					defkit.Field("mountPath", defkit.ParamTypeString).Required(),
					defkit.Field("subPath", defkit.ParamTypeString),
					defkit.Field("defaultMode", defkit.ParamTypeInt).Default(420),
					defkit.Field("secretName", defkit.ParamTypeString).Required(),
					defkit.Field("items", defkit.ParamTypeArray).
						ArrayOf(defkit.ParamTypeStruct).
						Nested(defkit.Struct("").Fields(
							defkit.Field("key", defkit.ParamTypeString).Required(),
							defkit.Field("path", defkit.ParamTypeString).Required(),
							defkit.Field("mode", defkit.ParamTypeInt).Default(511),
						)),
				)).
				Description("Mount Secret type storage"),
			defkit.Field("hostPath", defkit.ParamTypeArray).
				ArrayOf(defkit.ParamTypeStruct).
				Nested(defkit.Struct("").Fields(
					defkit.Field("name", defkit.ParamTypeString).Required(),
					defkit.Field("path", defkit.ParamTypeString).Required(),
					defkit.Field("mountPath", defkit.ParamTypeString).Required(),
					defkit.Field("type", defkit.ParamTypeString).
						Default("Directory").
						Enum("Directory", "DirectoryOrCreate", "FileOrCreate", "File", "Socket", "CharDevice", "BlockDevice"),
				)).
				Description("Declare host path type storage"),
		)

	// Build names used in templates
	jobName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName, defkit.Lit("-"), stepSessionID)
	containerName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName, defkit.Lit("-"), stepSessionID, defkit.Lit("-job"))
	stepLabel := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName)

	// Volume mount helpers using Reference for comprehension-based CUE patterns
	mountsArray := defkit.Reference(`[
	if parameter.storage != _|_ && parameter.storage.secret != _|_ for v in parameter.storage.secret {
		{
			name:      "secret-" + v.name
			mountPath: v.mountPath
			if v.subPath != _|_ {
				subPath: v.subPath
			}
		}
	},
	if parameter.storage != _|_ && parameter.storage.hostPath != _|_ for v in parameter.storage.hostPath {
		{
			name:      "hostpath-" + v.name
			mountPath: v.mountPath
		}
	},
]`)

	volumesList := defkit.Reference(`[
	if parameter.storage != _|_ && parameter.storage.secret != _|_ for v in parameter.storage.secret {
		{
			name: "secret-" + v.name
			secret: {
				defaultMode: v.defaultMode
				secretName:  v.secretName
				if v.items != _|_ {
					items: v.items
				}
			}
		}
	},
	if parameter.storage != _|_ && parameter.storage.hostPath != _|_ for v in parameter.storage.hostPath {
		{
			name: "hostpath-" + v.name
			path: v.path
		}
	},
]`)

	deDupVolumesArray := defkit.Reference(`[
	for val in [
		for i, vi in volumesList {
			for j, vj in volumesList if j < i && vi.name == vj.name {
				_ignore: true
			}
			vi
		},
	] if val._ignore == _|_ {
		val
	},
]`)

	// Job resource value
	jobValue := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("batch/v1")).
		Set("kind", defkit.Lit("Job")).
		Set("metadata", defkit.NewArrayElement().
			Set("name", jobName).
			SetIf(serviceAccountName.Eq("kubevela-vela-core"), "namespace", defkit.Lit("vela-system")).
			SetIf(serviceAccountName.Ne("kubevela-vela-core"), "namespace", vela.Namespace()),
		).
		Set("spec", defkit.NewArrayElement().
			Set("backoffLimit", defkit.Lit(3)).
			Set("template", defkit.NewArrayElement().
				Set("metadata", defkit.NewArrayElement().
					Set("labels", defkit.NewArrayElement().
						Set(`"workflow.oam.dev/step-name"`, stepLabel),
					),
				).
				Set("spec", defkit.NewArrayElement().
					Set("containers", defkit.NewArray().Item(
						defkit.NewArrayElement().
							Set("name", containerName).
							Set("image", image).
							Set("command", command).
							Set("volumeMounts", defkit.Reference("mountsArray")),
					)).
					Set("restartPolicy", defkit.Lit("Never")).
					Set("serviceAccount", serviceAccountName).
					Set("volumes", defkit.Reference("deDupVolumesArray")),
				),
			),
		)

	// Log resource selector
	logSelector := defkit.NewArrayElement().
		Set("source", defkit.NewArrayElement().
			Set("resources", defkit.NewArray().Item(
				defkit.NewArrayElement().
					Set("labelSelector", defkit.NewArrayElement().
						Set(`"workflow.oam.dev/step-name"`, stepLabel),
					),
			)),
		)

	return defkit.NewWorkflowStep("vela-cli").
		Description("Run a vela command").
		Category("Scripts & Commands").
		WithImports("vela/kube", "vela/builtin", "vela/util").
		Params(command, image, serviceAccountName, storage).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			// Volume computation helpers
			tpl.Set("mountsArray", mountsArray)
			tpl.Set("volumesList", volumesList)
			tpl.Set("deDupVolumesArray", deDupVolumesArray)

			// Apply the Job
			tpl.Builtin("job", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value": jobValue,
				}).
				Build()

			// Stream logs
			tpl.Builtin("log", "util.#Log").
				WithParams(map[string]defkit.Value{
					"source": logSelector,
				}).
				Build()

			// Fail if too many retries
			tpl.Set("fail", defkit.Reference(`{
	if job.$returns.value.status != _|_ if job.$returns.value.status.failed != _|_ {
		if job.$returns.value.status.failed > 2 {
			breakWorkflow: builtin.#Fail & {
				$params: message: "failed to execute vela command"
			}
		}
	}
}`))

			// Wait for Job success
			tpl.Set("wait", defkit.Reference(`builtin.#ConditionalWait & {
	if job.$returns.value.status != _|_ if job.$returns.value.status.succeeded != _|_ {
		$params: continue: job.$returns.value.status.succeeded > 0
	}
}`))
		})
}

func init() {
	defkit.Register(VelaCli())
}
