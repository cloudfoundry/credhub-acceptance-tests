package integration

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"os"
)

var _ = Describe("Sets debug flag", func() {

	BeforeEach(func() {
		os.Setenv("CREDHUB_DEBUG", "true")
	})

	It("should print debug info", func() {

		By("generating a credential", func() {
			name := GenerateUniqueCredentialName()
			session := RunCommand("generate", "-n", name, "-t", "user")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`[DEBUG]`))
		})

		By("setting a credential", func() {
			name := GenerateUniqueCredentialName()
			session := RunCommand("set", "-n", name, "-t", "password", "-w", "value")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`[DEBUG]`))
		})
	})

	AfterEach(func() {
		os.Unsetenv("CREDHUB_DEBUG")
	})

})
