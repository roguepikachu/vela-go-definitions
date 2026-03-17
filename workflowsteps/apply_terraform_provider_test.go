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

package workflowsteps_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("ApplyTerraformProvider WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.ApplyTerraformProvider()
		Expect(step.GetName()).To(Equal("apply-terraform-provider"))
		Expect(step.GetDescription()).To(Equal("Apply terraform provider config"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyTerraformProvider()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, alias, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Terraform"`))
			Expect(cueOutput).To(ContainSubstring(`alias: ""`))
			Expect(cueOutput).To(ContainSubstring(`"apply-terraform-provider": {`))
		})

		It("should import vela/config, vela/kube, and vela/builtin", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
		})

		It("should define all 8 provider helpers with default names and type constraints", func() {
			Expect(cueOutput).To(ContainSubstring("#AlibabaProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#AWSProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#AzureProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#BaiduProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#ECProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#GCPProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#TencentProvider: {"))
			Expect(cueOutput).To(ContainSubstring("#UCloudProvider: {"))

			Expect(cueOutput).To(ContainSubstring(`name: *"alibaba-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"aws-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"azure-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"baidu-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"ec-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"gcp-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"tencent-provider" | string`))
			Expect(cueOutput).To(ContainSubstring(`name: *"ucloud-provider" | string`))

			Expect(cueOutput).To(ContainSubstring(`type: "alibaba"`))
			Expect(cueOutput).To(ContainSubstring(`type: "aws"`))
			Expect(cueOutput).To(ContainSubstring(`type: "baidu"`))
			Expect(cueOutput).To(ContainSubstring(`type: "ec"`))
			Expect(cueOutput).To(ContainSubstring(`type: "gcp"`))
			Expect(cueOutput).To(ContainSubstring(`type: "tencent"`))
			Expect(cueOutput).To(ContainSubstring(`type: "ucloud"`))
		})

		It("should mark accessKey, secretKey, region as required in providers with providerBasic", func() {
			alibabaIdx := strings.Index(cueOutput, "#AlibabaProvider: {")
			Expect(alibabaIdx).To(BeNumerically(">", 0))
			alibabaBlock := cueOutput[alibabaIdx : alibabaIdx+200]
			Expect(alibabaBlock).To(ContainSubstring("accessKey!: string"))
			Expect(alibabaBlock).To(ContainSubstring("secretKey!: string"))
			Expect(alibabaBlock).To(ContainSubstring("region!: string"))

			awsIdx := strings.Index(cueOutput, "#AWSProvider: {")
			Expect(awsIdx).To(BeNumerically(">", 0))
			awsBlock := cueOutput[awsIdx : awsIdx+200]
			Expect(awsBlock).To(ContainSubstring("accessKey!: string"))
			Expect(awsBlock).To(ContainSubstring("secretKey!: string"))
			Expect(awsBlock).To(ContainSubstring("region!: string"))

			baiduIdx := strings.Index(cueOutput, "#BaiduProvider: {")
			Expect(baiduIdx).To(BeNumerically(">", 0))
			baiduBlock := cueOutput[baiduIdx : baiduIdx+200]
			Expect(baiduBlock).To(ContainSubstring("accessKey!: string"))
			Expect(baiduBlock).To(ContainSubstring("secretKey!: string"))
			Expect(baiduBlock).To(ContainSubstring("region!: string"))
		})

		It("should declare parameter as a union of all provider helpers", func() {
			Expect(cueOutput).To(ContainSubstring("parameter: #AlibabaProvider | #AWSProvider | #AzureProvider | #BaiduProvider | #ECProvider | #GCPProvider | #TencentProvider | #UCloudProvider"))
		})

		It("should create config and conditionally set provider-specific keys", func() {
			Expect(cueOutput).To(ContainSubstring("config.#CreateConfig & {"))
			Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepName)`))
			Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			Expect(cueOutput).To(ContainSubstring(`"terraform-\(parameter.type)"`))

			Expect(cueOutput).To(ContainSubstring(`parameter.type == "alibaba"`))
			Expect(cueOutput).To(ContainSubstring("ALICLOUD_ACCESS_KEY: parameter.accessKey"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "aws"`))
			Expect(cueOutput).To(ContainSubstring("AWS_ACCESS_KEY_ID: parameter.accessKey"))
			Expect(cueOutput).To(ContainSubstring("AWS_SESSION_TOKEN: parameter.token"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "azure"`))
			Expect(cueOutput).To(ContainSubstring("ARM_CLIENT_ID: parameter.clientID"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "gcp"`))
			Expect(cueOutput).To(ContainSubstring("GOOGLE_CREDENTIALS: parameter.credentials"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "baidu"`))
			Expect(cueOutput).To(ContainSubstring("BAIDUCLOUD_ACCESS_KEY: parameter.accessKey"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "ec"`))
			Expect(cueOutput).To(ContainSubstring("EC_API_KEY: parameter.apiKey"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "tencent"`))
			Expect(cueOutput).To(ContainSubstring("TENCENTCLOUD_SECRET_ID: parameter.secretID"))
			Expect(cueOutput).To(ContainSubstring(`parameter.type == "ucloud"`))
			Expect(cueOutput).To(ContainSubstring("UCLOUD_PRIVATE_KEY: parameter.privateKey"))
		})

		It("should read terraform Provider and wait for ready state", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "terraform.core.oam.dev/v1beta1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Provider"`))
			Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			Expect(cueOutput).To(ContainSubstring("read.$returns.value.status != _|_"))
			Expect(cueOutput).To(ContainSubstring(`read.$returns.value.status.state == "ready"`))
		})

		It("should have exactly one of each action type", func() {
			Expect(strings.Count(cueOutput, "config.#CreateConfig & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "kube.#Read & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "builtin.#ConditionalWait & {")).To(Equal(1))
		})
	})
})
