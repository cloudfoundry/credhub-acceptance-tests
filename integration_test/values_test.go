package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ =	It("should set, get, and delete a new value secret", func() {
	credentialName := generateUniqueCredentialName()

	By("trying to access a secret that doesn't exist", func() {
		session := runCommand("get", "-n", credentialName)
		stdErr := string(session.Err.Contents())

		Expect(stdErr).To(MatchRegexp(`Credential not found. Please validate your input and retry your request.`))
		Eventually(session).Should(Exit(1))
	})

	By("setting a new value secret", func() {
		session := runCommand("set", "-n", credentialName, "-t", "value", "-v", credentialValue)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
		Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
	})

	By("getting the new value secret", func() {
		session := runCommand("get", "-n", credentialName)
		stdOut := string(session.Out.Contents())

		Eventually(session).Should(Exit(0))

		Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
		Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
	})

	By("deleting the secret", func() {
		session := runCommand("delete", "-n", credentialName)
		Eventually(session).Should(Exit(0))
	})
})
