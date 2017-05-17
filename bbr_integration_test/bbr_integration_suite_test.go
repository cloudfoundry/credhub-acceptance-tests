package bbr_integration

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"io/ioutil"
	"testing"
	"time"

	"os/exec"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

func TestBbrIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup and Restore integration suite")
}

var config test_helpers.Config
var tmpDir string

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)

	var err error
	config, err = test_helpers.LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	tmpDir, err = ioutil.TempDir("", "BBR_CREDHUB_TEST")
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	Expect(os.RemoveAll(tmpDir)).To(Succeed())
})

func RunCommand(args ...string) *gexec.Session {
	fmt.Printf("Running %s", args)
	cmd := exec.Command(args[0], args[1:]...)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}
