package commands_test

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
	"regexp"
	"fmt"
)

var (
	commandPath string
	homeDir     string
	cfg         Config
	err         error
)

var _ = Describe("Integration test", func() {
	BeforeEach(func() {
		cfg, err = loadConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	It("smoke tests ok", func() {
		session := runCommand("api", cfg.ApiUrl)
		Eventually(session).Should(Exit(0))

		session = runCommand("login", "-u", "credhub_cli", "-p", "credhub_cli_password")
		Eventually(session).Should(Exit(0))

		uniqueId := strconv.FormatInt(time.Now().UnixNano(), 10)

		session = runCommand("get", "-n", uniqueId)
		Expect(session.Err.Contents()).To(MatchRegexp(`Secret not found. Please validate your input and retry your request.`))
		Eventually(session).Should(Exit(1))

		session = runCommand("set", "-n", uniqueId, "-t", "value", "-v", "bar")
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+value`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Value:\s+bar`))

		r, _ := regexp.Compile(`Updated:\s+(.*)[\s|$]`)
		original_timestamp_array := r.FindSubmatch(session.Out.Contents())

		Expect(original_timestamp_array).To(HaveLen(2))

		original_timestamp := original_timestamp_array[1]

		Expect(original_timestamp).NotTo(HaveLen(0))

		session = runCommand("set", "-n", uniqueId, "-t", "value", "-v", "newvalue", "--no-overwrite")
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+value`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Value:\s+bar`))
		Expect(session.Out.Contents()).To(MatchRegexp(fmt.Sprintf(`Updated:\s+%s`, original_timestamp)))

		// We need to sleep in order to ensure that the timestamp is different,
		// since it is truncated to the second.
		time.Sleep(time.Duration(1) * time.Second)

		session = runCommand("set", "-n", uniqueId, "-t", "value", "-v", "newvalue")
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+value`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Value:\s+newvalue`))
		Expect(session.Out.Contents()).NotTo(MatchRegexp(fmt.Sprintf(`Updated:\s+%s`, original_timestamp)))

		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+value`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Value:\s+newvalue`))

		uniqueCertificateId := strconv.FormatInt(time.Now().UnixNano(), 10)

		session = runCommand("set", "-n", uniqueCertificateId, "-t", "certificate", "--certificate-string", "iamacertificate")
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+certificate`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Certificate:\s+iamacertificate`))

		session = runCommand("set", "-n", uniqueCertificateId, "-t", "certificate", "--no-overwrite")
		Eventually(session).Should(Exit(1))
		Expect(session.Err.Contents()).To(MatchRegexp(".*At least one certificate type must be set. Please validate your input and retry your request."))

		session = runCommand("delete", "-n", uniqueId)
		Eventually(session).Should(Exit(0))

		uniqueId2 := uniqueId + "2"
		session = runCommand("get", "-n", uniqueId2)
		Expect(session.Err.Contents()).To(MatchRegexp(`Secret not found. Please validate your input and retry your request.`))
		Eventually(session).Should(Exit(1))

		session = runCommand("ca-get", "-n", uniqueId)
		Expect(session.Err.Contents()).To(MatchRegexp(`CA not found. Please validate your input and retry your request.`))
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

		session = runCommand("get", "-n", uniqueId2)
		Eventually(session).Should(Exit(0))

		uniqueId3 := uniqueId + "3"
		session = runCommand("generate", "-n", uniqueId3, "-t", "rsa")
		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(MatchRegexp(`Type:\s+rsa`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Public Key:\s+-----BEGIN PUBLIC KEY-----`))
		Expect(session.Out.Contents()).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))

		session = runCommand("get", "-n", uniqueId3)
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
	ApiUrl string `json:"api_url"`
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
