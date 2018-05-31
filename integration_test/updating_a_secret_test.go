package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("updating a secret", func() {
	Describe("updating with a set (PUT)", func() {
		It("should be able to overwrite a secret", func() {
			credentialName := GenerateUniqueCredentialName()

			By("setting a new value secret", func() {
				RunCommand("set", "-n", credentialName, "-t", "value", "-v", "old value")
				session := RunCommand("get", "-n", credentialName)
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(`type: value`))
				Expect(stdOut).To(ContainSubstring("value: old value"))
			})

			By("setting the value secret again", func() {
				RunCommand("set", "-n", credentialName, "-t", "value", "-v", "new value")
				session := RunCommand("get", "-n", credentialName)
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(`type: value`))
				Expect(stdOut).To(ContainSubstring("value: new value"))
			})
		})
	})

	Describe("generating -> setting -> generating", func() {
		It("does not bleed values from the generate", func() {
			caName := GenerateUniqueCredentialName()
			credentialname := GenerateUniqueCredentialName()

			By("generating a new ca", func() {
				RunCommand("generate", "-n", caName, "-t", "certificate", "-c", "anything", "--is-ca", "--self-sign")
			})

			By("generating a new certificate signed by the CA", func() {
				RunCommand("generate", "-n", credentialname, "-t", "certificate", "-c", "bla", "--ca", caName)
				session := RunCommand("get", "-n", credentialname)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(MatchRegexp(`certificate:\ |\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`ca: |\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
			})

			By("overwriting the certificate with `set`", func() {
				RunCommand("set", "-n", credentialname, "-t", "certificate", "--certificate", VALID_CERTIFICATE)
				session := RunCommand("get", "-n", credentialname)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(ContainSubstring(`certificate: `))
				Expect(stdOut).To(ContainSubstring(VALID_CERTIFICATE_OUTPUT))
			})
		})
	})
})
