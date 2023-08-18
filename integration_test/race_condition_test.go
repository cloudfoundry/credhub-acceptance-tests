package integration_test

import (
	"strings"
	"time"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Race condition tests", func() {
	Describe("when generating a new secret in multiple threads with `--no-overwrite`", func() {
		It("should return the same value for both", func() {
			rsaSecretName := GenerateUniqueCredentialName()

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				time.Sleep(150 * time.Millisecond)
				RunCommand("generate", "-n", rsaSecretName, "-t", "rsa", "--no-overwrite")
				session := RunCommand("get", "-n", rsaSecretName)
				waitForSession1 <- session
			}()

			go func() {
				RunCommand("generate", "-n", rsaSecretName, "-t", "rsa", "--no-overwrite")
				session := RunCommand("get", "-n", rsaSecretName)
				waitForSession2 <- session
			}()

			session1 := <-waitForSession1
			session2 := <-waitForSession2

			Eventually(session1).Should(Exit(0))
			Eventually(session2).Should(Exit(0))
			stdOut1 := string(session1.Out.Contents())
			stdOut2 := string(session2.Out.Contents())

			Expect(stdOut1).To(ContainSubstring(`type: rsa`))
			Expect(stdOut1).To(MatchRegexp(`public_key: |\s+-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut1).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))

			Expect(stdOut1).To(Equal(stdOut2))
		})
	})

	Describe("when setting a new secret in multiple threads", func() {
		It("should return the different values for both", func() {
			passwordSecretName := GenerateUniqueCredentialName()

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				session := RunCommand("set", "-n", passwordSecretName, "-w", "test-value", "-t", "password")
				waitForSession1 <- session
			}()

			go func() {
				session := RunCommand("set", "-n", passwordSecretName, "-w", "test-value", "-t", "password")
				waitForSession2 <- session
			}()

			session1 := <-waitForSession1
			session2 := <-waitForSession2

			Eventually(session1).Should(Exit(0))
			Eventually(session2).Should(Exit(0))
			stdOut1 := string(session1.Out.Contents())
			stdOut2 := string(session2.Out.Contents())

			Expect(stdOut1).NotTo(Equal(stdOut2))
		})
	})

	Describe("when setting one secret name for two types", func() {
		It("should return a type mismatch error", func() {
			rsaSecretName := GenerateUniqueCredentialName()
			type_error := "The credential type cannot be modified. Please delete the credential if you wish to create it with a different type."

			waitForSession1 := make(chan *Session)
			waitForSession2 := make(chan *Session)

			go func() {
				session := RunCommand("set", "-n", rsaSecretName, "-t", "ssh", "-p", "something")
				waitForSession1 <- session
			}()

			go func() {
				session := RunCommand("set", "-n", rsaSecretName, "-t", "rsa", "-p", "something")
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
