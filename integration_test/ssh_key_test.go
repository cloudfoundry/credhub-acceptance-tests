package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("SSH key test", func() {
	Describe("setting an SSH key", func() {
		It("should be able to set an ssh key", func() {
			base64DecodablePublicKey := "public"
			session := RunCommand("set", "-n", GenerateUniqueCredentialName(), "-t", "ssh", "-U", base64DecodablePublicKey, "-P", credentialValue)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+` + base64DecodablePublicKey))
			Expect(stdOut).To(MatchRegexp("Private Key:\\s+" + credentialValue))
		})
	})

	It("should generate an SSH key", func() {
		sshSecretName := GenerateUniqueCredentialName()

		By("generating the key", func() {
			session := RunCommand("generate", "-n", sshSecretName, "-t", "ssh", "-m", "some comment")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+ssh-rsa \S+`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := RunCommand("get", "-n", sshSecretName)
			Eventually(session).Should(Exit(0))
		})
	})
})
