package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Password test", func() {
	It("should generate a password", func() {
		session := runCommand("generate", "-n", generateUniqueCredentialName(), "-t", "password")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
	})

	It("should regenerate passwords with similar rules", func() {
		generatedPasswordId := generateUniqueCredentialName()

		By("first generating a password with no numbers", func() {
			session := runCommand("generate", "-n", generatedPasswordId, "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(Not(MatchRegexp(`Value:\s.+\d`)))
		})

		By("then regenerating the password and observing it still has no numbers", func() {
			session := runCommand("regenerate", "-n", generatedPasswordId)
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(Not(MatchRegexp(`Value:\s.+\d`)))
		})
	})
})
