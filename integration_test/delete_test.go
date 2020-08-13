package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Delete Test", func() {
	It("should delete credentials by name", func() {
		RunCommand("generate", "-n", "/a/cred1", "-t", "password")

		session := RunCommand("delete", "-n", "/a/cred1")
		Eventually(session).Should(Exit(0))

		session = RunCommand("get", "-n", "/a/cred1")
		Eventually(session).Should(Exit(1))
		Expect(string(session.Err.Contents())).To(ContainSubstring("The request could not be completed because the credential does not exist or you do not have sufficient authorization."))
	})

	It("should delete credentials by path", func() {
		RunCommand("generate", "-n", "/a1", "-t", "password")
		RunCommand("generate", "-n", "/a1/cred1", "-t", "password")
		RunCommand("generate", "-n", "/a1/b/cred2", "-t", "password")
		RunCommand("generate", "-n", "/a1/b/cred3", "-t", "password")

		session := RunCommand("delete", "-p", "/a1")
		Eventually(session).Should(Exit(0))

		session = RunCommand("find", "-p", "/a1")
		Eventually(session).Should(Exit(0))
		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring("credentials: []"))

		session = RunCommand("get", "-n", "/a1")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring("/a1"))
	})
})
