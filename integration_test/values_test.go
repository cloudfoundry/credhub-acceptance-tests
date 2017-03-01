package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ =	It("should set, get, and delete a new value secret", func() {
	credentialName := GenerateUniqueCredentialName()

	By("trying to access a secret that doesn't exist", func() {
		session := RunCommand("get", "-n", credentialName)
		stdErr := string(session.Err.Contents())

		Expect(stdErr).To(MatchRegexp(`Credential not found. Please validate your input and retry your request.`))
		Eventually(session).Should(Exit(1))
	})

	By("setting a new value secret", func() {
		session := RunCommand("set", "-n", credentialName, "-t", "value", "-v", credentialValue)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
		Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
	})

	By("getting the new value secret", func() {
		session := RunCommand("get", "-n", credentialName)
		stdOut := string(session.Out.Contents())

		Eventually(session).Should(Exit(0))

		Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
		Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
	})

	By("deleting the secret", func() {
		session := RunCommand("delete", "-n", credentialName)
		Eventually(session).Should(Exit(0))
	})
})
