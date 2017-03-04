package integration

import (
	"os/exec"
	"strings"
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
				session := runCommandWithMTLS(config, config.ValidPKCS12Path)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(MatchRegexp(`"type":"password"`))
			})
		})
	})
})

func runCommandWithMTLS(config Config, pkcs12Path string) *Session {
	url := strings.Join([]string{config.ApiUrl, "/api/v1/data"}, "")
	payload := `{"name":"mtlstest","type":"password"}`
	println("pkcs12Path", pkcs12Path)
	content_type := "Content-Type: application/json"
	cmd := exec.Command("/usr/local/opt/curl/bin/curl",  "-k", url, "-H", content_type, "-X", "POST", "-d", payload, "--cert", pkcs12Path)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}
