package remote_backend

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"strings"
)

var _ = Describe("Find", func() {
	Describe("containing name", func() {
		BeforeEach(func() {
			session := RunCommand("set", "-t", "value", "-n", "/some-credential", "-v", "value")
			Expect(session).Should(Exit(0))

			session = RunCommand("set", "-t", "value", "-n", "/some-other-credential", "-v", "value")
			Expect(session).Should(Exit(0))

			session = RunCommand("set", "-t", "value", "-n", "/another-credential", "-v", "value")
			Expect(session).Should(Exit(0))

		})
		Context("when credentials exist with matching name", func() {
			It("shows all credentials containing name", func() {
				session := RunCommand("find", "-n", "other")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).ToNot(ContainSubstring("- name: /some-credential"))
				Expect(stdOut).To(ContainSubstring("- name: /some-other-credential"))
				Expect(stdOut).To(ContainSubstring("- name: /another-credential"))

			})
		})
		Context("when no credentials exist with matching name", func() {
			It("returns error message", func() {
				session := RunCommand("find", "-n", "abc")
				Expect(session).Should(Exit(1))
				stdOut := string(session.Err.Contents())
				Expect(stdOut).To(ContainSubstring("No credentials exist which match the provided parameters."))
			})
		})

	})
	Describe("starting with path", func() {
		BeforeEach(func() {
			session := RunCommand("set", "-t", "value", "-n", "/some/credential", "-v", "value")
			Expect(session).Should(Exit(0))

			session = RunCommand("set", "-t", "value", "-n", "/some/other-credential", "-v", "value")
			Expect(session).Should(Exit(0))

			session = RunCommand("set", "-t", "value", "-n", "/another/credential", "-v", "value")
			Expect(session).Should(Exit(0))

		})
		Context("when credentials exist starting with path", func() {
			It("shows all credentials starting with path", func() {
				session := RunCommand("find", "-p", "/some")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring("- name: /some/credential"))
				Expect(stdOut).To(ContainSubstring("- name: /some/other-credential"))
				Expect(stdOut).ToNot(ContainSubstring("- name: /another/credential"))
			})
		})
		Context("when no credentials exist starting with path", func() {
			It("returns error message", func() {
				session := RunCommand("find", "-p", "/abc")
				Expect(session).Should(Exit(0))
				stdOut := strings.TrimSpace(string(session.Out.Contents()))
				Expect(stdOut).To(ContainSubstring("credentials: []"))
			})
		})
	})
})
