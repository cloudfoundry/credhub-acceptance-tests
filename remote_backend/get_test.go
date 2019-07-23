package remote_backend_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Get", func() {
	Describe("get by name", func() {
		It("gets a cred by name", func() {
			session := RunCommand("get", "-n", "my-value")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("my-value"))
			Expect(stdOut).To(ContainSubstring("test-value"))
		})
		Context("getting all versions", func() {
			It("gets all versions of the credential", func() {
				credentialName := GenerateUniqueCredentialName()
				session := RunCommand("set", "-n", credentialName, "-t", "password", "--password", "first-password")
				Eventually(session).Should(Exit(0))
				session = RunCommand("set", "-n", credentialName, "-t", "password", "--password", "second-password")
				Eventually(session).Should(Exit(0))

				session = RunCommand("curl", "-p", fmt.Sprintf("/api/v1/data?name=%s", credentialName))
				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(`"value": "second-password"`))
				Expect(stdOut).To(ContainSubstring(`"value": "first-password"`))
			})
		})
		Context("getting multiple versions", func() {
			It("gets the desired number of versions of the credential", func() {
				credentialName := GenerateUniqueCredentialName()
				session := RunCommand("set", "-n", credentialName, "-t", "password", "--password", "first-password")
				Eventually(session).Should(Exit(0))
				session = RunCommand("set", "-n", credentialName, "-t", "password", "--password", "second-password")
				Eventually(session).Should(Exit(0))

				session = RunCommand("curl", "-p", fmt.Sprintf("/api/v1/data?name=%s&versions=1", credentialName))
				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(`"value": "second-password"`))
				Expect(stdOut).NotTo(ContainSubstring(`"value": "first-password"`))
			})
		})
	})
	Describe("get by id", func() {
		It("gets a cred by id", func() {
			session := RunCommand("get", "--id", "00000000-0000-0000-0000-000000000001")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("my-value"))
			Expect(stdOut).To(ContainSubstring("test-value"))
		})
	})
})
