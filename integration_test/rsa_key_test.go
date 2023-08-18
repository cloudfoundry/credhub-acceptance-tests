package integration_test

import (
	"strings"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("RSA key test", func() {
	Describe("setting an RSA key", func() {
		It("should be able to set an rsa key", func() {
			credName := GenerateUniqueCredentialName()
			RunCommand("set", "-n", credName, "-t", "rsa", "-u", "iamapublickey", "-p", credentialValue)
			session := RunCommand("get", "-n", credName)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`type: rsa`))
			Expect(stdOut).To(ContainSubstring(`public_key: iamapublickey`))
			Expect(stdOut).To(ContainSubstring("private_key: " + credentialValue))
		})
	})

	It("should generate an RSA key", func() {
		rsaSecretName := GenerateUniqueCredentialName()

		By("generating the key", func() {
			session := RunCommand("generate", "-n", rsaSecretName, "-t", "rsa")
			Eventually(session).Should(Exit(0))

			session = RunCommand("get", "-n", rsaSecretName)
			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(ContainSubstring(`type: rsa`))
			Expect(stdOut).To(MatchRegexp(`public_key: |\s+-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := RunCommand("get", "-n", rsaSecretName)
			Eventually(session).Should(Exit(0))
		})
	})

	It("should regenerate an RSA key", func() {
		rsaSecretName := GenerateUniqueCredentialName()

		By("regenerate should create an new value", func() {
			RunCommand("generate", "-n", rsaSecretName, "-t", "rsa")
			session := RunCommand("get", "-n", rsaSecretName)

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())
			initialPublicKey := stdOut[strings.Index(stdOut, "-----BEGIN PUBLIC KEY-----"):strings.Index(stdOut, "-----END PUBLIC KEY-----")]
			initialPrivateKey := stdOut[strings.Index(stdOut, "-----BEGIN RSA PRIVATE KEY-----"):strings.Index(stdOut, "-----END RSA PRIVATE KEY-----")]
			session = RunCommand("regenerate", "-n", rsaSecretName)
			Eventually(session).Should(Exit(0))

			session = RunCommand("get", "-n", rsaSecretName)
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).NotTo(ContainSubstring(initialPublicKey))
			Expect(stdOut).NotTo(ContainSubstring(initialPrivateKey))
		})
	})
})
