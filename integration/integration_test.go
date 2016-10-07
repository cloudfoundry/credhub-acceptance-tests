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

		// These happen before each test due to the lack of a BeforeAll
		// (https://github.com/onsi/ginkgo/issues/70) :(
		// If the tests are slow, they should be runnable in parallel with the -p option.
		session := runCommand("api", cfg.ApiUrl)
		Eventually(session).Should(Exit(0))

		session = runCommand("login", "-u", "credhub_cli", "-p", "credhub_cli_password")
		Eventually(session).Should(Exit(0))
	})

	It("should set, get, and delete a new value secret", func() {
		valueId := strconv.FormatInt(time.Now().UnixNano(), 10)

		By("trying to access a secret that doesn't exist", func() {
			session := runCommand("get", "-n", valueId)
			stdErr := string(session.Err.Contents())

			Expect(stdErr).To(MatchRegexp(`Secret not found. Please validate your input and retry your request.`))
			Eventually(session).Should(Exit(1))
		})

		By("setting a new value secret", func() {
			session := runCommand("set", "-n", valueId, "-t", "value", "-v", "bar")
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp(`Value:\s+bar`))
		})

		By("getting the new value secret", func() {
			session := runCommand("get", "-n", valueId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp(`Value:\s+bar`))
		})

		By("deleting the secret", func() {
			session := runCommand("delete", "-n", valueId)
			Eventually(session).Should(Exit(0))
		})
	})

	It("should set a secret's timestamp correctly", func() {
		var original_timestamp []byte
		r, _ := regexp.Compile(`Updated:\s+(.*)[\s|$]`)

		secretId := strconv.FormatInt(time.Now().UnixNano(), 10)

		By("getting the original timestamp", func() {
			session := runCommand("set", "-n", secretId, "-t", "value", "-v", "bar")
			original_timestamp_array := r.FindSubmatch(session.Out.Contents())

			Expect(original_timestamp_array).To(HaveLen(2))

			original_timestamp = original_timestamp_array[1]

			Expect(original_timestamp).NotTo(HaveLen(0))
		})

		By("getting the timestamp after a no-overwrite set", func() {
			session := runCommand("set", "-n", secretId, "-t", "value", "-v", "newvalue", "--no-overwrite")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp(`Value:\s+bar`))
			Expect(stdOut).To(MatchRegexp(fmt.Sprintf(`Updated:\s+%s`, original_timestamp)))
		})

		By("getting the timestamp after an overwrite set", func() {
			// We need to sleep in order to ensure that the timestamp is different,
			// since it is truncated to the second.
			time.Sleep(time.Duration(1) * time.Second)

			session := runCommand("set", "-n", secretId, "-t", "value", "-v", "newvalue")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp(`Value:\s+newvalue`))
			Expect(stdOut).NotTo(MatchRegexp(fmt.Sprintf(`Updated:\s+%s`, original_timestamp)))
		})

		By("getting the value", func() {
			session := runCommand("get", "-n", secretId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp(`Value:\s+newvalue`))
		})
	})

	Describe("setting a certificate", func() {
		certificateId := strconv.FormatInt(time.Now().UnixNano(), 10)

		It("should be able to set a certificate", func() {
			session := runCommand("set", "-n", certificateId, "-t", "certificate", "--certificate-string", "iamacertificate")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
			Expect(stdOut).To(MatchRegexp(`Certificate:\s+iamacertificate`))
		})

		It("should require a certificate type", func() {
			session := runCommand("set", "-n", certificateId, "-t", "certificate", "--no-overwrite")
			Eventually(session).Should(Exit(1))
			Expect(session.Err.Contents()).To(MatchRegexp(".*At least one certificate type must be set. Please validate your input and retry your request."))
		})
	})

	It("should generate a CA and certificate", func() {
		certificateAuthorityId := strconv.FormatInt(time.Now().UnixNano(), 10)
		certificateId := certificateAuthorityId + "1"

		By("retrieving a CA that doesn't exist yet", func() {
			session := runCommand("ca-get", "-n", certificateAuthorityId)
			stdErr := string(session.Err.Contents())

			Expect(stdErr).To(MatchRegexp(`CA not found. Please validate your input and retry your request.`))
			Eventually(session).Should(Exit(1))
		})

		By("generating the CA", func() {
			session := runCommand("ca-generate", "-n", certificateAuthorityId, "--common-name", certificateAuthorityId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+root`))
			Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
		})

		By("getting the new CA", func() {
			session := runCommand("ca-get", "-n", certificateAuthorityId)
			Eventually(session).Should(Exit(0))
		})

		By("generating the certificate", func() {
			session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
			Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
		})

		By("getting the certificate", func() {
			session := runCommand("get", "-n", certificateId)
			Eventually(session).Should(Exit(0))
		})
	})

	It("should generate an SSH key", func() {
		sshId := strconv.FormatInt(time.Now().UnixNano(), 10)

		By("generating the key", func() {
			session := runCommand("generate", "-n", sshId, "-t", "ssh")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+ssh-rsa \S+`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := runCommand("get", "-n", sshId)
			Eventually(session).Should(Exit(0))
		})
	})

	It("should generate an RSA key", func() {
		rsaId := strconv.FormatInt(time.Now().UnixNano(), 10)

		By("generating the key", func() {
			session := runCommand("generate", "-n", rsaId, "-t", "rsa")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+rsa`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := runCommand("get", "-n", rsaId)
			Eventually(session).Should(Exit(0))
		})
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
