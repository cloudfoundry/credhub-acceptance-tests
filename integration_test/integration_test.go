package integration_test

import (
	"strconv"
	"time"

	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"encoding/json"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	commandPath string
	homeDir     string
	cfg         Config
)

// We look for these values in the verify-logging CI task to ensure that credentials don't leak
const credentialValue = "FAKE-CREDENTIAL-VALUE"

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
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

	cfg, err = loadConfig()
	Expect(err).NotTo(HaveOccurred())

	// These happen before each test due to the lack of a BeforeAll
	// (https://github.com/onsi/ginkgo/issues/70) :(
	// If the tests are slow, they should be runnable in parallel with the -p option.
	session := runCommand("api", cfg.ApiUrl)
	Eventually(session).Should(Exit(0))

	session = runCommand("login", "-u", cfg.ApiUsername, "-p", cfg.ApiPassword)
	Eventually(session).Should(Exit(0))
})

var _ = AfterEach(func() {
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/pivotal-cf/credhub-cli")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	commandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})

func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

type Config struct {
	ApiUrl      string `json:"api_url"`
	ApiUsername string `json:"api_username"`
	ApiPassword string `json:"api_password"`
}

func generateUniqueCredentialName() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func loadConfig() (Config, error) {
	c := Config{}

	data, err := ioutil.ReadFile(path.Join(os.Getenv("PWD"), "config.json"))
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
