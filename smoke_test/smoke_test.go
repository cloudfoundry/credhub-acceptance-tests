package smoke_test
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Smoke Test", func() {

	Describe("certificates", func() {
		certificate := "smoke_test_value" + GenerateUniqueCredentialName()
		It("can CRD certificates", func() {
			By("should be able to set a certificate", func() {
				session := RunCommand("set", "-n", certificate, "-t", "certificate", "--certificate-string", "iamacertificate")
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(MatchRegexp(`certificate: |\s+iamacertificate`))
			})

			By("should be able to get the certificate", func() {
				session := RunCommand("get", "-n", certificate)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(MatchRegexp(`certificate: |\s+iamacertificate`))
			})

			By("should be able to delete the certificate", func() {
				session := RunCommand("delete", "-n", certificate)

				Eventually(session).Should(Exit(0))

				session = RunCommand("get", "-n", certificate)
				stdErr := string(session.Err.Contents())

				Eventually(session).Should(Exit(1))

				Expect(stdErr).To(ContainSubstring(`Credential not found. Please validate your input and retry your request.`))
			})
		})
	})
})

