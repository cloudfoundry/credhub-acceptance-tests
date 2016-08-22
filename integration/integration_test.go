package commands_test

import (
	"strconv"
	"time"

	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
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

		uniqueId := strconv.FormatInt(time.Now().UnixNano(), 10)

		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(1))

		session = runCommand("set", "-n", uniqueId, "-v", "bar")
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+value`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Value:\s+bar`))

		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(0))

		session = runCommand("delete", "-n", uniqueId)
		Eventually(session).Should(Exit(0))

		uniqueId2 := uniqueId + "2"
		session = runCommand("get", "-n", uniqueId2)
		Eventually(session).Should(Exit(1))

		session = runCommand("ca-get", "-n", uniqueId)
		Eventually(session).Should(Exit(1))

		session = runCommand("ca-generate", "-n", uniqueId, "--common-name", uniqueId)
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+root`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))

		session = runCommand("ca-get", "-n", uniqueId)
		Eventually(session).Should(Exit(0))

		session = runCommand("generate", "-n", uniqueId2, "-t", "certificate", "--common-name", uniqueId2, "--ca", uniqueId)
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+certificate`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))

		session = runCommand("get", "-n", uniqueId2)
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
