package integration_test

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"testing"

	"github.com/hashicorp/go-version"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	homeDir string
	cfg     Config
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

	os.Unsetenv("CREDHUB_DEBUG")

	cfg, err = LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	// These happen before each test due to the lack of a BeforeAll
	// (https://github.com/onsi/ginkgo/issues/70) :(
	// If the tests are slow, they should be runnable in parallel with the -p option.
	TargetAndLoginWithClientCredentials(cfg)
})

var _ = AfterEach(func() {
	CleanEnv()
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("code.cloudfoundry.org/credhub-cli", "-mod=mod")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	CommandPath = string(data)

	rand.Seed(GinkgoRandomSeed() + int64(GinkgoParallelNode()))
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})

func getServerVersion() (*version.Version, error) {
	session := RunCommand("--version")
	if session.ExitCode() != 0 {
		return nil, errors.New(string(session.Err.Contents()))
	}

	r := regexp.MustCompile(`Server Version: (\d+\.\d+\.\d+)`)
	matches := r.FindSubmatch(session.Out.Contents())
	if len(matches) < 2 {
		return nil, errors.New("failed to find semver version for credhub server")
	}
	v, err := version.NewVersion(string(matches[1]))
	if err != nil {
		return nil, err
	}

	return v, nil
}
