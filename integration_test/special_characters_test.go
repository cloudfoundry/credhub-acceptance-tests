package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("handling special characters", func() {
	It("should handle secrets whose names begin with a leading slash", func() {
		baseId := "ace/ventura" + GenerateUniqueCredentialName()
		leadingSlashId := "/" + baseId
		passwordValue := "finkel-is-einhorn"

		By("setting a value whose name begins with a leading slash", func() {
			session := RunCommand("set", "-n", leadingSlashId, "-t", "password", "-w", passwordValue)
			Eventually(session).Should(Exit(0))
		})

		By("retrieving the value that was set with a leading slash", func() {
			session := RunCommand("get", "-n", leadingSlashId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`type: password`))
			Expect(stdOut).To(ContainSubstring(passwordValue))
		})

		By("retrieving the value that was set without a leading slash", func() {
			session := RunCommand("get", "-n", baseId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`type: password`))
			Expect(stdOut).To(ContainSubstring(passwordValue))
		})
	})
})
