package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("updating a secret", func() {
	Describe("updating with a set (PUT)", func() {
		It("should be able to overwrite a secret", func() {
			credentialName := generateUniqueCredentialName()

			By("setting a new value secret", func() {
				session := runCommand("set", "-n", credentialName, "-t", "value", "-v", "old value")
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
				Expect(stdOut).To(MatchRegexp("Value:\\s+" + "old value"))
			})

			By("setting the value secret again", func() {
				session := runCommand("set", "-n", credentialName, "-t", "value", "-v", "new value")
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
				Expect(stdOut).To(MatchRegexp("Value:\\s+" + "new value"))
			})
		})
	})
})
