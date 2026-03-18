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

var _ = Describe("GenerateJDBCConnection WorkflowStep", func() {
	It("should have the correct name and description", func() {
		step := workflowsteps.GenerateJDBCConnection()
		Expect(step.GetName()).To(Equal("generate-jdbc-connection"))
		Expect(step.GetDescription()).To(Equal("Generate a JDBC connection based on Component of alibaba-rds"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.GenerateJDBCConnection()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct step header with type, category, and quoted name", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "Terraform"`))
			Expect(cueOutput).To(ContainSubstring(`"generate-jdbc-connection": {`))
		})

		It("should import vela/kube, vela/util, and encoding/base64", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
		})

		It("should declare name and optional namespace parameters", func() {
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("namespace?: string"))
		})

		It("should read a v1 Secret via kube.#Read with conditional namespace", func() {
			Expect(cueOutput).To(ContainSubstring("kube.#Read & {"))
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind: "Secret"`))
			Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			Expect(cueOutput).To(ContainSubstring("parameter.namespace != _|_"))
			Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
		})

		It("should decode all 5 DB fields via util.#ConvertString with base64.Decode", func() {
			Expect(cueOutput).To(ContainSubstring("dbHost: util.#ConvertString & {"))
			Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_HOST"]`))
			Expect(cueOutput).To(ContainSubstring("dbPort: util.#ConvertString & {"))
			Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_PORT"]`))
			Expect(cueOutput).To(ContainSubstring("dbName: util.#ConvertString & {"))
			Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_NAME"]`))
			Expect(cueOutput).To(ContainSubstring("username: util.#ConvertString & {"))
			Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_USER"]`))
			Expect(cueOutput).To(ContainSubstring("password: util.#ConvertString & {"))
			Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_PASSWORD"]`))
			Expect(strings.Count(cueOutput, "base64.Decode(null,")).To(Equal(5))
		})

		It("should build JDBC URL env array with decoded parts", func() {
			Expect(cueOutput).To(ContainSubstring("env: ["))
			Expect(cueOutput).To(ContainSubstring(`"jdbc://" + dbHost.$returns.str`))
			Expect(cueOutput).To(ContainSubstring(`dbPort.$returns.str`))
			Expect(cueOutput).To(ContainSubstring(`dbName.$returns.str`))
			Expect(cueOutput).To(ContainSubstring("characterEncoding=utf8&useSSL=false"))
			Expect(cueOutput).To(ContainSubstring(`name: "username", value: username.$returns.str`))
			Expect(cueOutput).To(ContainSubstring(`name: "password", value: password.$returns.str`))
		})

		It("should be structurally correct", func() {
			Expect(strings.Count(cueOutput, "kube.#Read & {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "util.#ConvertString & {")).To(Equal(5))
		})
	})
})
