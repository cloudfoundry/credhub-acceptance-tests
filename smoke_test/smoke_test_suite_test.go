package smoke_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"io/ioutil"
	"runtime"
	"os"
)

var (
	homeDir     string
	cfg         Config
)

var _ = BeforeEach(func() {
	CommandPath = "credhub"

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

	TargetAndLoginSkipTls(cfg)
})

func TestSmokeTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SmokeTest Suite")
}

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})
