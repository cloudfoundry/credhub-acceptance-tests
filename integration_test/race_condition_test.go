package integration_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Race condition tests", func() {
	Describe("when creating a new secret in multiple threads with `--no-overwrite`", func() {
		It("should return the same value for both given they are seperated by at least 150 milliseconds", func() {
			rsaSecretName := generateUniqueCredentialName()

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				time.Sleep(150 * time.Millisecond)
				session := runCommand("generate", "-n", rsaSecretName, "-t", "rsa", "--no-overwrite")
				waitForSession1 <- session
			}()

			go func() {
				session := runCommand("generate", "-n", rsaSecretName, "-t", "rsa", "--no-overwrite")
				waitForSession2 <- session
			}()

			session1 := <-waitForSession1
			session2 := <-waitForSession2

			Eventually(session1).Should(Exit(0))
			Eventually(session2).Should(Exit(0))
			stdOut1 := string(session1.Out.Contents())
			stdOut2 := string(session2.Out.Contents())

			Expect(stdOut1).To(MatchRegexp(`Type:\s+rsa`))
			Expect(stdOut1).To(MatchRegexp(`Public Key:\s+-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut1).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))

			Expect(stdOut1).To(Equal(stdOut2))
		})
	})
})
