package integration_test

import (
	"time"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/pivotal-cf/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Race condition tests", func() {
	Describe("when generating a new secret in multiple threads with `--no-overwrite`", func() {
		It("should return the same value for both", func() {
			rsaSecretName := GenerateUniqueCredentialName()

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				time.Sleep(150 * time.Millisecond)
				session := RunCommand("generate", "-n", rsaSecretName, "-t", "rsa", "--no-overwrite")
				waitForSession1 <- session
			}()

			go func() {
				session := RunCommand("generate", "-n", rsaSecretName, "-t", "rsa", "--no-overwrite")
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

	Describe("when setting a new secret in multiple threads with `--no-overwrite`", func() {
		It("should return the same value for both", func() {
			rsaSecretName := GenerateUniqueCredentialName()

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				session := RunCommand("set", "-n", rsaSecretName, "-v", "test-value", "--no-overwrite")
				waitForSession1 <- session
			}()

			go func() {
				session := RunCommand("set", "-n", rsaSecretName, "-v", "test-value", "--no-overwrite")
				waitForSession2 <- session
			}()

			session1 := <-waitForSession1
			session2 := <-waitForSession2

			Eventually(session1).Should(Exit(0))
			Eventually(session2).Should(Exit(0))
			stdOut1 := string(session1.Out.Contents())
			stdOut2 := string(session2.Out.Contents())

			Expect(stdOut1).To(Equal(stdOut2))
		})
	})

	Describe("when setting one secret name for two types", func() {
		It("should return a type mismatch error", func() {
			rsaSecretName := GenerateUniqueCredentialName()
			type_error := "The credential type cannot be modified. Please delete the credential if you wish to create it with a different type."

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				session := RunCommand("set", "-n", rsaSecretName, "-t", "ssh", "-P", "something")
				waitForSession1 <- session
			}()

			go func() {
				session := RunCommand("set", "-n", rsaSecretName, "-t", "rsa", "-P", "something")
				waitForSession2 <- session
			}()

			session1 := <-waitForSession1
			session2 := <-waitForSession2

			Eventually(session1).Should(Exit())
			Eventually(session2).Should(Exit())
			out1 := string(session1.Out.Contents()) + string(session1.Err.Contents())
			out2 := string(session2.Out.Contents()) + string(session2.Err.Contents())
			
			errors_out := strings.Contains(out1, type_error) || strings.Contains(out2, type_error)
			Expect(errors_out).To(BeTrue())
		})
	})
})
