package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"strings"
)

var _ = Describe("SSH key test", func() {
	Describe("setting an SSH key", func() {
		It("should be able to set an ssh key", func() {
			base64DecodablePublicKey := "public"
			credName := GenerateUniqueCredentialName()
			RunCommand("set", "-n", credName, "-t", "ssh", "-u", base64DecodablePublicKey, "-p", credentialValue)
			session := RunCommand("get", "-n", credName)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`type: ssh`))
			Expect(stdOut).To(ContainSubstring(`public_key: ` + base64DecodablePublicKey))
			Expect(stdOut).To(ContainSubstring("private_key: " + credentialValue))
		})
	})

	It("should generate an SSH key", func() {
		sshSecretName := GenerateUniqueCredentialName()

		By("generating the key", func() {
			RunCommand("generate", "-n", sshSecretName, "-t", "ssh", "-m", "some comment")
			session := RunCommand("get", "-n", sshSecretName)

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(ContainSubstring(`type: ssh`))
			Expect(stdOut).To(MatchRegexp(`public_key: ssh-rsa \w+`))
			Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := RunCommand("get", "-n", sshSecretName)
			Eventually(session).Should(Exit(0))
		})
	})

	It("should regenerate an SSH key", func() {
		sshSecretName := GenerateUniqueCredentialName()

		By("regenerate should create a new value", func() {
			RunCommand("generate", "-n", sshSecretName, "-t", "ssh", "-m", "some comment")
			session := RunCommand("get", "-n", sshSecretName)

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())
			initialPublicKey := stdOut[strings.Index(stdOut, "public_key: ssh-rsa"):strings.Index(stdOut, "some comment")]
			initialPrivateKey := stdOut[strings.Index(stdOut, "-----BEGIN RSA PRIVATE KEY-----"):strings.Index(stdOut, "-----END RSA PRIVATE KEY-----")]

			RunCommand("regenerate", "-n", sshSecretName)
			session = RunCommand("get", "-n", sshSecretName)

			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`type: ssh`))
			Expect(stdOut).To(MatchRegexp(`public_key: ssh-rsa \w+`))
			Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
			Expect(stdOut).NotTo(ContainSubstring(initialPublicKey))
			Expect(stdOut).NotTo(ContainSubstring(initialPrivateKey))
		})

	})
})
