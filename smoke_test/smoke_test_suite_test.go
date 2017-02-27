package smoke_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/pivotal-cf/credhub-acceptance-tests/test_helpers"
	"io/ioutil"
	"runtime"
	"os"
)

var (
	homeDir     string
	cfg         Config
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
	session := RunCommand("api", cfg.ApiUrl)
	Eventually(session).Should(Exit(0))

	session = RunCommand("login", "-u", cfg.ApiUsername, "-p", cfg.ApiPassword)
	Eventually(session).Should(Exit(0))
})

func TestSmokeTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SmokeTest Suite")
}


var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/pivotal-cf/credhub-cli")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	CommandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})