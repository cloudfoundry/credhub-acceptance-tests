package test_helpers

import (
	"strconv"
	"time"

	"io/ioutil"
	"os"
	"os/exec"

	"encoding/json"
	"path"

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

type Config struct {
	ApiUrl      string `json:"api_url"`
	ApiUsername string `json:"api_username"`
	ApiPassword string `json:"api_password"`
}

func LoadConfig() (Config, error) {

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

func TargetAndLogin(cfg Config) {
	session := RunCommand("login", "-s", cfg.ApiUrl, "-u", cfg.ApiUsername, "-p", cfg.ApiPassword)
	Eventually(session).Should(Exit(0))
}
