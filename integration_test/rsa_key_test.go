package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("RSA key test", func() {
	Describe("setting an RSA key", func() {
		It("should be able to set an rsa key", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "rsa", "-U", "iamapublickey", "-P", credentialValue)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+rsa`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+iamapublickey`))
			Expect(stdOut).To(MatchRegexp("Private Key:\\s+" + credentialValue))
		})
	})

	It("should generate an RSA key", func() {
		rsaSecretName := generateUniqueCredentialName()

		By("generating the key", func() {
			session := runCommand("generate", "-n", rsaSecretName, "-t", "rsa")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+rsa`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := runCommand("get", "-n", rsaSecretName)
			Eventually(session).Should(Exit(0))
		})
	})
})
