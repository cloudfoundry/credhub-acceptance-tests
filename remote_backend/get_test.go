package remote_backend_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Get", func() {
	Describe("get by name", func() {
		It("gets a cred by name", func() {
			session := RunCommand("get", "-n", "my-value")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("my-value"))
			Expect(stdOut).To(ContainSubstring("test-value"))
		})
	})
	Describe("get by id", func() {
		It("gets a cred by id", func() {
			session := RunCommand("get", "--id", "00000000-0000-0000-0000-000000000001")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("my-value"))
			Expect(stdOut).To(ContainSubstring("test-value"))
		})
	})
})
