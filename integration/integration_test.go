package commands_test

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
)

var (
	commandPath string
	homeDir     string
)

var _ = Describe("Integration test", func() {
	It("smoke tests ok", func() {
		session := runCommand("api", "https://50.17.59.67:8844")
		Eventually(session).Should(Exit(0))

		session = runCommand("login", "-u", "credhub_cli", "-p", "credhub_cli_password")
		Eventually(session).Should(Exit(0))

		data, err := ioutil.ReadFile(path.Join(userHomeDir(), ".cm", "config.json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(ContainSubstring(`"AuthURL":"https://50.17.59.67:8443"`))

		uniqueId := strconv.FormatInt(time.Now().UnixNano(), 10)
		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(1))

		session = runCommand("set", "-n", uniqueId, "-v", "bar")
		Eventually(session).Should(Exit(0))

		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(0))

		fmt.Println(string(session.Out.Contents()))
		session = runCommand("delete", "-n", uniqueId)
		Eventually(session).Should(Exit(0))
	})
})

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
})

var _ = AfterEach(func() {
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/pivotal-cf/cm-cli")
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

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}
