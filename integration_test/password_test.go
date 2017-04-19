package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"regexp"
)

var _ = Describe("Password test", func() {
	It("should set a password", func() {
		session := RunCommand("set", "-n", GenerateUniqueCredentialName(), "-t", "password", "-v", "some_value")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: some_value`))
	})

	It("should generate a password", func() {
		session := RunCommand("generate", "-n", GenerateUniqueCredentialName(), "-t", "password")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
	})

	It("should regenerate passwords with similar rules", func() {
		generatedPasswordId := GenerateUniqueCredentialName()
		firstValue :=""
		valueRegexp := regexp.MustCompile(`value: \D*`)

		By("first generating a password with no numbers", func() {
			session := RunCommand("generate", "-n", generatedPasswordId, "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`type: password`))
			Expect(stdOut).NotTo(MatchRegexp(`value: \S*\d`))

			firstValue = valueRegexp.FindString(stdOut)
		})

		By("then regenerating the password and observing it still has no numbers", func() {
			session := RunCommand("regenerate", "-n", generatedPasswordId)
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).NotTo(MatchRegexp(`value: \S*\d`))
			Expect(stdOut).NotTo(ContainSubstring(firstValue))
		})
	})
})
