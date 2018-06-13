package integration_test

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var (
	homeDir     string
	cfg         Config
)

// We look for these values in the verify-logging CI task to ensure that credentials don't leak
const credentialValue = "FAKE-CREDENTIAL-VALUE"

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

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
	TargetAndLoginWithClientCredentials(cfg)
})

var _ = AfterEach(func() {
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("code.cloudfoundry.org/credhub-cli")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	CommandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})
