package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"os/exec"
)

var (
	config Config
	err    error
)

var _ = Describe("json secrets", func() {
	BeforeEach(func() {
		config, err = LoadConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should set and get a new json secret", func() {
		credentialName := GenerateUniqueCredentialName()

		By("setting a new json secret", func() {
			json := `{"type":"json","name":"` + credentialName + `","value":{"object":{"is":"complex"}}}`

			cmd := exec.Command("curl",
				"-k", config.ApiUrl + "/api/v1/data",
				"-H", "Content-Type: application/json",
				"-X", "PUT",
				"-d", json,
				"--cert", config.ValidCertPath,
				"--key", config.ValidPrivateKeyPath)

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)

			Expect(err).NotTo(HaveOccurred())
			<-session.Exited

			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`"type":\s*"json"`))
			Expect(stdOut).To(MatchRegexp(`"value":\s*{"object":{"is":"complex"}}`))
		})

		By("getting the new json secret", func() {
			cmd := exec.Command("curl",
				"-k", config.ApiUrl+"/api/v1/data?name="+credentialName,
				"-H", "Content-Type: application/json",
				"-XGET",
				"--cert", config.ValidCertPath,
				"--key", config.ValidPrivateKeyPath)

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)

			Expect(err).NotTo(HaveOccurred())
			<-session.Exited

			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`"type":\s*"json"`))
			Expect(stdOut).To(MatchRegexp(`"value":\s*{"object":{"is":"complex"}}`))
		})
	})
})
