package remote_backend_test

import (
	"math/rand"
	"testing"

	"io/ioutil"
	"os"
	"runtime"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	homeDir string
	cfg     Config
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

	TargetAndLoginSkipTls(cfg)
})

func TestRemoteBackendTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RemoteBackend Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	return []byte("credhub")
}, func(cli_path []byte) {
	CommandPath = string(cli_path)

	rand.Seed(GinkgoRandomSeed() + int64(GinkgoParallelNode()))
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})
