package bbr_integration

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"testing"
	"time"
)

func TestBbrIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup and Restore integration suite")
}

var credhubUrl string
var credhubApiUsername string
var credhubApiPassword string
var deploymentName string
var tmpDir string
var boshCertPath string
var bbrBinaryPath string
var credhubCliBinaryPath string

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)

	var err error
	tmpDir, err = ioutil.TempDir("", "BBR_CREDHUB_TEST")
	Expect(err).NotTo(HaveOccurred())

	credhubUrl = MustHaveEnv("CREDHUB_URL")
	credhubApiUsername = MustHaveEnv("CREDHUB_API_USERNAME")
	credhubApiPassword = MustHaveEnv("CREDHUB_API_PASSWORD")
	deploymentName = MustHaveEnv("CREDHUB_DEPLOYMENT_NAME")
	boshCertPath = MustHaveEnv("BOSH_CERT_PATH")
	bbrBinaryPath = MustHaveEnv("BBR_PATH")
	credhubCliBinaryPath = MustHaveEnv("CREDHUB_CLI_PATH")
})

var _ = AfterSuite(func() {
	Expect(os.RemoveAll(tmpDir)).To(Succeed())
})

func MustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	Expect(val).NotTo(BeEmpty(), "Need "+keyname+" for the test")
	return val
}
