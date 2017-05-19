package test_helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	CommandPath string
)

func GenerateUniqueCredentialName() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func RunCommand(args ...string) *Session {
	cmd := exec.Command(CommandPath, args...)

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

type BoshConfig struct {
	Host              string `json:"host"`
	SshUsername       string `json:"bosh_ssh_username"`
	SshPrivateKeyPath string `json:"bosh_ssh_private_key_path"`
}

type Config struct {
	Bosh           *BoshConfig `json:"bosh"`
	ApiUrl         string      `json:"api_url"`
	ApiUsername    string      `json:"api_username"`
	ApiPassword    string      `json:"api_password"`
	CredentialRoot string      `json:"credential_root"`
}

func LoadConfig() (Config, error) {

	configuration := Config{}

	configurationJson, err := ioutil.ReadFile(path.Join(os.Getenv("PWD"), "config.json"))
	if err != nil {
		return configuration, err
	}

	err = json.Unmarshal(configurationJson, &configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func TargetAndLogin(cfg Config) {
	session := RunCommand("login", "-s", cfg.ApiUrl, "-u", cfg.ApiUsername, "-p", cfg.ApiPassword, "--skip-tls-validation")
	Eventually(session).Should(Exit(0))
}
