package integration

import (
	"os/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"fmt"
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

	Describe("with a certificate signed by a trusted CA	", func() {
		Describe("when the certificate has a valid date range", func() {
			It("allows the user to hit an authenticated endpoint", func() {
				session := runCommandWithMTLS(
						config.ApiUrl + "/api/v1/data",
						config.ValidCertPath,
						config.ValidPrivateKeyPath)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(MatchRegexp(`"type":"password"`))
			})
		})
	})
})

func runCommandWithMTLS(url, certPath, keyPath string) *Session{

	Expect(certPath).NotTo(BeEmpty())
	Expect(keyPath).NotTo(BeEmpty())

	payload := `{"name":"mtlstest","type":"password"}`
	content_type := "Content-Type: application/json"
	cmd := exec.Command("curl",
		"-k", url,
		"-H", content_type,
		"-XPOST",
		"-d", payload,
		"--cert", certPath,
		"--key", keyPath)

	fmt.Printf("%#v\n", cmd.Args)

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}
