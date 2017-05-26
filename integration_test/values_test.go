package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"regexp"
)

var _ =	It("should set, get, and delete a new value secret", func() {
	credentialName := GenerateUniqueCredentialName()

	By("trying to access a secret that doesn't exist", func() {
		session := RunCommand("get", "-n", credentialName)
		stdErr := string(session.Err.Contents())

		Expect(stdErr).To(ContainSubstring(`Credential not found. Please validate your input and retry your request.`))
		Eventually(session).Should(Exit(1))
	})

	By("setting a new value secret", func() {
		session := RunCommand("set", "-n", credentialName, "-t", "value", "-v", credentialValue)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring("value: " + credentialValue))
	})

	By("getting the new value secret", func() {
		session := RunCommand("get", "-n", credentialName)
		stdOut := string(session.Out.Contents())

		Eventually(session).Should(Exit(0))

		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring("value: " + credentialValue))

		re := regexp.MustCompile("id: (.*?)\n")
		credentialId := re.FindStringSubmatch(stdOut)

		session = RunCommand("get", "--id", credentialId[1])
		stdOut = string(session.Out.Contents())

		Eventually(session).Should(Exit(0))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring("value: " + credentialValue))

	})

	By("deleting the secret", func() {
		session := RunCommand("delete", "-n", credentialName)
		Eventually(session).Should(Exit(0))
	})
})
