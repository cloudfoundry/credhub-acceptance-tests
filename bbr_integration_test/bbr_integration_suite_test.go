package bbr_integration

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"testing"
	"time"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

func TestBbrIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup and Restore integration suite")
}

var config test_helpers.Config
var tmpDir string
var bbrBinaryPath string
var credhubCliBinaryPath string

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)

	var err error
	config, err = test_helpers.LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	tmpDir, err = ioutil.TempDir("", "BBR_CREDHUB_TEST")
	Expect(err).NotTo(HaveOccurred())

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
