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

var _ = Describe("Notification WorkflowStep", func() {
	It("should have correct metadata", func() {
		step := workflowsteps.Notification()
		Expect(step.GetName()).To(Equal("notification"))
		Expect(step.GetDescription()).To(ContainSubstring("Send notifications to Email, DingTalk, Slack, Lark"))
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Notification()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate step header with type and category", func() {
			Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cueOutput).To(ContainSubstring(`"category": "External Integration"`))
		})

		It("should import all required packages", func() {
			Expect(cueOutput).To(ContainSubstring(`"vela/http"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/email"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
			Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
		})

		It("should define all helper types with correct fields and cross-references", func() {
			// #TextType
			Expect(cueOutput).To(ContainSubstring("#TextType: {"))
			Expect(cueOutput).To(ContainSubstring("type: string"))
			Expect(cueOutput).To(ContainSubstring("text: string"))
			Expect(cueOutput).To(ContainSubstring("emoji?: bool"))
			Expect(cueOutput).To(ContainSubstring("verbatim?: bool"))

			// #Option referencing #TextType
			Expect(cueOutput).To(ContainSubstring("#Option: {"))
			Expect(cueOutput).To(ContainSubstring("text: #TextType"))
			Expect(cueOutput).To(ContainSubstring("description?: #TextType"))

			// #DingLink
			Expect(cueOutput).To(ContainSubstring("#DingLink: {"))
			Expect(cueOutput).To(ContainSubstring("messageUrl?: string"))
			Expect(cueOutput).To(ContainSubstring("picUrl?: string"))

			// #DingBtn
			Expect(cueOutput).To(ContainSubstring("#DingBtn: {"))
			Expect(cueOutput).To(ContainSubstring("title: string"))
			Expect(cueOutput).To(ContainSubstring("actionURL: string"))

			// #Block with element references to #TextType and #Option
			Expect(cueOutput).To(ContainSubstring("#Block: {"))
			Expect(cueOutput).To(ContainSubstring("block_id?: string"))
			Expect(cueOutput).To(ContainSubstring("elements?: [...{"))
			Expect(cueOutput).To(ContainSubstring("text?: #TextType"))
			Expect(cueOutput).To(ContainSubstring("placeholder?: #TextType"))
			Expect(cueOutput).To(ContainSubstring("options?: [...#Option]"))
			Expect(cueOutput).To(ContainSubstring("initial_options?: [...#Option]"))
			Expect(cueOutput).To(ContainSubstring("option_groups?: [...#Option]"))
		})

		It("should define lark parameter with url ClosedUnion and message fields", func() {
			Expect(cueOutput).To(ContainSubstring("lark?: {"))

			larkIdx := strings.Index(cueOutput, "lark?: {")
			Expect(larkIdx).To(BeNumerically(">", 0))
			larkBlock := cueOutput[larkIdx:]
			Expect(larkBlock).To(ContainSubstring("url: close({"))
			Expect(larkBlock).To(ContainSubstring("value: string"))
			Expect(larkBlock).To(ContainSubstring("}) | close({"))
			Expect(larkBlock).To(ContainSubstring("secretRef: {"))

			Expect(cueOutput).To(ContainSubstring("msg_type: string"))
			Expect(cueOutput).To(ContainSubstring("// +usage=content should be json encode string"))
		})

		It("should define dingding parameter with all message types and helpers", func() {
			Expect(cueOutput).To(ContainSubstring("dingding?: {"))

			// url as ClosedUnion
			dingIdx := strings.Index(cueOutput, "dingding?: {")
			Expect(dingIdx).To(BeNumerically(">", 0))
			dingBlock := cueOutput[dingIdx:]
			Expect(dingBlock).To(ContainSubstring("url: close({"))
			Expect(dingBlock).To(ContainSubstring("value: string"))
			Expect(dingBlock).To(ContainSubstring("}) | close({"))

			// message with msgtype enum
			Expect(cueOutput).To(ContainSubstring(`msgtype: *"text" | "link" | "markdown" | "actionCard" | "feedCard"`))

			// text, link, markdown, at, actionCard, feedCard
			Expect(cueOutput).To(ContainSubstring("text?: close({"))
			Expect(cueOutput).To(ContainSubstring("content: string"))
			Expect(cueOutput).To(ContainSubstring("link?: #DingLink"))
			Expect(cueOutput).To(ContainSubstring("markdown?: close({"))
			Expect(cueOutput).To(ContainSubstring("at?: close({"))
			Expect(cueOutput).To(ContainSubstring("atMobiles?: [...string]"))
			Expect(cueOutput).To(ContainSubstring("isAtAll?: bool"))
			Expect(cueOutput).To(ContainSubstring("actionCard?: close({"))
			Expect(cueOutput).To(ContainSubstring("hideAvatar: string"))
			Expect(cueOutput).To(ContainSubstring("btnOrientation: string"))
			Expect(cueOutput).To(ContainSubstring("singleTitle: string"))
			Expect(cueOutput).To(ContainSubstring("singleURL: string"))
			Expect(cueOutput).To(ContainSubstring("#DingBtn"))
			Expect(cueOutput).To(ContainSubstring("feedCard?: close({"))
			Expect(cueOutput).To(ContainSubstring("links: [...#DingLink]"))
		})

		It("should define slack parameter with message, blocks, attachments, and options", func() {
			Expect(cueOutput).To(ContainSubstring("slack?: {"))
			Expect(cueOutput).To(ContainSubstring("// +usage=Specify the message text for slack notification"))
			Expect(cueOutput).To(ContainSubstring("blocks?: [...#Block]"))
			Expect(cueOutput).To(ContainSubstring("attachments?: close({"))
			Expect(cueOutput).To(ContainSubstring("color?: string"))
			Expect(cueOutput).To(ContainSubstring("thread_ts?: string"))
			Expect(cueOutput).To(ContainSubstring("mrkdwn?: *true | bool"))
		})

		It("should define email parameter with from, password ClosedUnion, to, and content", func() {
			Expect(cueOutput).To(ContainSubstring("email?: {"))

			// from fields
			Expect(cueOutput).To(ContainSubstring("address: string"))
			Expect(cueOutput).To(ContainSubstring("alias?: string"))
			Expect(cueOutput).To(ContainSubstring("host: string"))
			Expect(cueOutput).To(ContainSubstring("port: *587 | int"))

			// password as ClosedUnion
			Expect(cueOutput).To(ContainSubstring("// +usage=Specify the password of the email"))
			emailIdx := strings.Index(cueOutput, "email?: {")
			Expect(emailIdx).To(BeNumerically(">", 0))
			emailBlock := cueOutput[emailIdx:]
			Expect(emailBlock).To(ContainSubstring("password: close({"))
			Expect(emailBlock).To(ContainSubstring("}) | close({"))

			// to and content
			Expect(cueOutput).To(ContainSubstring("to: [...string]"))
			Expect(cueOutput).To(ContainSubstring("subject: string"))
			Expect(cueOutput).To(ContainSubstring("body: string"))
		})

		It("should generate guarded channel blocks with guards inside field scope", func() {
			Expect(cueOutput).To(ContainSubstring("ding: {"))
			Expect(cueOutput).To(ContainSubstring("if parameter.dingding != _|_ {"))
			Expect(cueOutput).To(ContainSubstring("lark: {"))
			Expect(cueOutput).To(ContainSubstring("if parameter.lark != _|_ {"))
			Expect(cueOutput).To(ContainSubstring("slack: {"))
			Expect(cueOutput).To(ContainSubstring("if parameter.slack != _|_ {"))
			Expect(cueOutput).To(ContainSubstring("email0: {"))
			Expect(cueOutput).To(ContainSubstring("if parameter.email != _|_ {"))

			// Verify guard is inside field, not outside
			dingFieldIdx := strings.Index(cueOutput, "ding: {")
			dingGuardIdx := strings.Index(cueOutput, "if parameter.dingding != _|_")
			Expect(dingFieldIdx).To(BeNumerically(">", 0))
			Expect(dingGuardIdx).To(BeNumerically(">", dingFieldIdx))
		})

		It("should generate dingding channel template actions with value and secretRef paths", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.dingding.url.value != _|_"))
			Expect(cueOutput).To(ContainSubstring("url:    parameter.dingding.url.value"))
			Expect(cueOutput).To(ContainSubstring("parameter.dingding.url.secretRef != _|_ && parameter.dingding.url.value == _|_"))
			Expect(cueOutput).To(ContainSubstring("name:      parameter.dingding.url.secretRef.name"))
			Expect(cueOutput).To(ContainSubstring("base64.Decode(null, read.$returns.value.data[parameter.dingding.url.secretRef.key])"))
			Expect(cueOutput).To(ContainSubstring("url:    stringValue.$returns.str"))
			Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.dingding.message)"))
		})

		It("should generate lark channel template actions with value and secretRef paths", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.lark.url.value != _|_"))
			Expect(cueOutput).To(ContainSubstring("url:    parameter.lark.url.value"))
			Expect(cueOutput).To(ContainSubstring("parameter.lark.url.secretRef != _|_ && parameter.lark.url.value == _|_"))
			Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.lark.message)"))
		})

		It("should generate slack channel template actions with value and secretRef paths", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.slack.url.value != _|_"))
			Expect(cueOutput).To(ContainSubstring("url:    parameter.slack.url.value"))
			Expect(cueOutput).To(ContainSubstring("parameter.slack.url.secretRef != _|_ && parameter.slack.url.value == _|_"))
			Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.slack.message)"))
		})

		It("should generate email template actions with password value and secretRef paths", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.email.from.password.value != _|_"))
			Expect(cueOutput).To(ContainSubstring("email.#SendEmail"))
			Expect(cueOutput).To(ContainSubstring("parameter.email.from.password.secretRef != _|_ && parameter.email.from.password.value == _|_"))
			Expect(cueOutput).To(ContainSubstring("address: parameter.email.from.address"))
			Expect(cueOutput).To(ContainSubstring("host:     parameter.email.from.host"))
			Expect(cueOutput).To(ContainSubstring("if parameter.email.from.alias != _|_"))
			Expect(cueOutput).To(ContainSubstring("alias: parameter.email.from.alias"))
			Expect(cueOutput).To(ContainSubstring("to:      parameter.email.to"))
			Expect(cueOutput).To(ContainSubstring("content: parameter.email.content"))
		})

		It("should have correct structural counts for operations across all channels", func() {
			// 2 HTTP ops per channel (value URL + secretRef URL) x 3 channels (ding, lark, slack)
			Expect(strings.Count(cueOutput, "http.#HTTPDo & {")).To(Equal(6))

			// 1 kube.#Read per channel (ding, lark, slack) + 1 for email = 4
			Expect(strings.Count(cueOutput, "kube.#Read & {")).To(Equal(4))

			// 1 ConvertString per channel (ding, lark, slack) + 1 for email = 4
			Expect(strings.Count(cueOutput, "util.#ConvertString & {")).To(Equal(4))

			// 2 email.#SendEmail: one for password.value, one for password secretRef
			Expect(strings.Count(cueOutput, "email.#SendEmail")).To(Equal(2))

			// Content-Type header on all 6 HTTP posts
			Expect(strings.Count(cueOutput, `header: "Content-Type": "application/json"`)).To(Equal(6))

			// Exactly 4 guarded blocks: ding, lark, slack, email0
			Expect(strings.Count(cueOutput, "ding: {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "lark: {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "slack: {")).To(Equal(1))
			Expect(strings.Count(cueOutput, "email0: {")).To(Equal(1))
		})
	})
})
