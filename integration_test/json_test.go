package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("json secrets", func() {
	credentialName := GenerateUniqueCredentialName()
	credentialValue := `{"object":{"is":"complex"},"has":["an","array"]}`
	credentialYaml := "value:\n    has:\n        - an\n        - array\n    object:\n        is: complex\n"

	It("should set, get, and delete a new json secret", func() {
		By("setting a new json secret", func() {
			RunCommand("set", "-n", credentialName, "-t", "json", "-v", credentialValue)
			session := RunCommand("get", "-n", credentialName)
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`type: json`))
			Expect(stdOut).To(ContainSubstring(credentialYaml))
		})

		By("getting the new json secret", func() {
			session := RunCommand("get", "-n", credentialName)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`type: json`))
			Expect(stdOut).To(ContainSubstring(credentialYaml))
		})

		By("deleting the secret", func() {
			session := RunCommand("delete", "-n", credentialName)
			Eventually(session).Should(Exit(0))
		})
	})
})
