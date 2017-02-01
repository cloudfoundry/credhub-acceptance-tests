package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("handling special characters", func() {
	It("should handle secrets who names begin with a leading slash", func() {
		baseId := "ace/ventura" + generateUniqueCredentialName()
		leadingSlashId := "/" + baseId
		passwordValue := "finkel-is-einhorn"

		By("setting a value whose name begins with a leading slash", func() {
			session := runCommand("set", "-n", leadingSlashId, "-t", "password", "-v", passwordValue)
			Eventually(session).Should(Exit(0))
		})

		By("retrieving the value that was set with a leading slash", func() {
			session := runCommand("get", "-n", leadingSlashId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(ContainSubstring(passwordValue))
		})

		By("retrieving the value that was set without a leading slash", func() {
			session := runCommand("get", "-n", baseId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(ContainSubstring(passwordValue))
		})
	})

	It("should get secrets whose names have lots of special characters", func() {
		crazyCharsId := "dan:test/ing?danother[stuff]that@shouldn!tbe$in&the" + generateUniqueCredentialName()

		By("setting a value with lots of special characters", func() {
			session := runCommand("set", "-n", crazyCharsId, "-t", "password", "-v", "woof-woof")
			Eventually(session).Should(Exit(0))
		})

		By("retrieving the value that was set", func() {
			session := runCommand("get", "-n", crazyCharsId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(ContainSubstring(crazyCharsId))
		})
	})

	It("should handle edge-casey character combinations", func() {
		edgeCaseId := "&gunk=x/bar/cr@zytown108" + generateUniqueCredentialName()

		By("setting a value with lots of special characters", func() {
			session := runCommand("set", "-n", edgeCaseId, "-t", "password", "-v", "find-me")
			Eventually(session).Should(Exit(0))
		})

		By("retrieving the value that was set", func() {
			session := runCommand("get", "-n", edgeCaseId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(ContainSubstring(edgeCaseId))
		})
	})

	It("should delete secrets with special characters", func() {
		deleteId := "?testParam=foo&gunk=x/bar/piv0t@l" + generateUniqueCredentialName()

		By("setting a value with lots of special characters", func() {
			session := runCommand("set", "-n", deleteId, "-t", "password", "-v", "find-me")
			Eventually(session).Should(Exit(0))
		})

		By("deleting the secret", func() {
			session := runCommand("delete", "-n", deleteId)
			Eventually(session).Should(Exit(0))
		})
	})
})
