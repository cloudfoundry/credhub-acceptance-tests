package integration

import (
	"os/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("mutual TLS authentication", func() {
	var (
		config Config
		err error
	)

	BeforeEach(func() {
		config, err = LoadConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("with a certicate signed by a trusted CA	", func() {
		Describe("when the certificate has a valid date range", func() {
			It("allows the user to hit an authenticated endpoint", func() {
				session := runCommandWithMTLS(config)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(MatchRegexp(`"type":"password"`))
			})
		})
	})
})

func runCommandWithMTLS(config Config) *Session {
	url := config.ApiUrl + "/api/v1/data"
	pemPathWithPassword := config.ValidPEMPath + ":" + config.MTLSPassword
	payload := `{"name":"mtlstest","type":"password"}`
	println("pem certificate path", pemPathWithPassword)
	content_type := "Content-Type: application/json"
	cmd := exec.Command("curl",  "-k", url, "-H", content_type, "-X", "POST", "-d", payload, "--cert", pemPathWithPassword)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}
