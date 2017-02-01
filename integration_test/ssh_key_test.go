package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("SSH key test", func() {
	Describe("setting an SSH key", func() {
		It("should be able to set an ssh key", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "ssh", "-U", "iamapublickey", "-P", credentialValue)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+iamapublickey`))
			Expect(stdOut).To(MatchRegexp("Private Key:\\s+" + credentialValue))
		})
	})

	It("should generate an SSH key", func() {
		sshSecretName := generateUniqueCredentialName()

		By("generating the key", func() {
			session := runCommand("generate", "-n", sshSecretName, "-t", "ssh")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+ssh-rsa \S+`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := runCommand("get", "-n", sshSecretName)
			Eventually(session).Should(Exit(0))
		})
	})
})
