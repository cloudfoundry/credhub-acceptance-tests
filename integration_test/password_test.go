package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Password test", func() {
	It("should set a password", func() {
		session := RunCommand("set", "-n", GenerateUniqueCredentialName(), "-t", "password", "-v", "some_value")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
		Expect(stdOut).To(MatchRegexp(`Value:\s+some_value`))
	})

	It("should generate a password", func() {
		session := RunCommand("generate", "-n", GenerateUniqueCredentialName(), "-t", "password")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
	})

	It("should regenerate passwords with similar rules", func() {
		generatedPasswordId := GenerateUniqueCredentialName()

		By("first generating a password with no numbers", func() {
			session := RunCommand("generate", "-n", generatedPasswordId, "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(Not(MatchRegexp(`Value:\s.+\d`)))
		})

		By("then regenerating the password and observing it still has no numbers", func() {
			session := RunCommand("regenerate", "-n", generatedPasswordId)
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(Not(MatchRegexp(`Value:\s.+\d`)))
		})
	})
})
