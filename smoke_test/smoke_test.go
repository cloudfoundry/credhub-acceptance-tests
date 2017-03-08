package smoke_test
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"io/ioutil"
	"runtime"
	"os"
)

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "cm-test")
	Expect(err).NotTo(HaveOccurred())

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", homeDir)
	} else {
		os.Setenv("HOME", homeDir)
	}

	cfg, err = LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	// These happen before each test due to the lack of a BeforeAll
	// (https://github.com/onsi/ginkgo/issues/70) :(
	// If the tests are slow, they should be runnable in parallel with the -p option.
	TargetAndLogin(cfg)
})

var _ = Describe("Smoke Test", func() {

	Describe("certificates", func() {
		certificate := "smoke_test_value" + GenerateUniqueCredentialName()
		It("can CRD certificates", func() {
			By("should be able to set a certificate", func() {
				session := RunCommand("set", "-n", certificate, "-t", "certificate", "--certificate-string", "iamacertificate")
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+iamacertificate`))
			})

			By("should be able to get the certificate", func() {
				session := RunCommand("get", "-n", certificate)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+iamacertificate`))
			})

			By("should be able to delete the certificate", func() {
				session := RunCommand("delete", "-n", certificate)

				Eventually(session).Should(Exit(0))

				session = RunCommand("get", "-n", certificate)
				stdErr := string(session.Err.Contents())

				Eventually(session).Should(Exit(1))

				Expect(stdErr).To(MatchRegexp(`Credential not found. Please validate your input and retry your request.`))
			})
		})
	})
})

