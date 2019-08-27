package test_helpers

import (
  "encoding/json"
  "io/ioutil"
  "math/rand"
  "os"
  "os/exec"
  "path"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/onsi/gomega/gexec"
)

var (
  CommandPath string
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateUniqueCredentialName() string {
  length := 20
  name := make([]byte, length)
  for i := range name {
    name[i] = chars[rand.Intn(len(chars))]
  }
  return string(name)
}

func RunCommand(args ...string) *Session {
  cmd := exec.Command(CommandPath, args...)

  session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
  Expect(err).NotTo(HaveOccurred())
  <-session.Exited

  return session
}

type BoshConfig struct {
  Environment  string `json:"bosh_environment"`
  Client       string `json:"bosh_client"`
  ClientSecret string `json:"bosh_client_secret"`
  CaCertPath   string `json:"bosh_ca_cert_path"`
}

type Config struct {
  Bosh           *BoshConfig `json:"bosh"`
  ApiUrl         string      `json:"api_url"`
  ApiUsername    string      `json:"api_username"`
  ApiPassword    string      `json:"api_password"`
  CredentialRoot string      `json:"credential_root"`
  UAACa          string      `json:"uaa_ca"`
  DirectorHost   string      `json:"director_host"`
  ClientName     string      `json:"client_name"`
  ClientSecret   string      `json:"client_secret"`
  DeploymentName string      `json:"deployment_name"`
  ConcatenateCas bool        `json:"concatenate_cas"`
}

func LoadConfig() (Config, error) {

  configuration := Config{}

  configurationJson, err := ioutil.ReadFile(path.Join(os.Getenv("PWD"), "test_config.json"))
  if err != nil {
    return configuration, err
  }

  err = json.Unmarshal(configurationJson, &configuration)
  if err != nil {
    return configuration, err
  }

  return configuration, nil
}

func TargetAndLoginWithClientCredentials(cfg Config) {
  CleanEnv()
  credhub_ca := path.Join(cfg.CredentialRoot, "server_ca_cert.pem")
  uaa_ca := cfg.UAACa
  credhub_ca_contents, _ := ioutil.ReadFile(credhub_ca)
  uaa_ca_contents, _ := ioutil.ReadFile(uaa_ca)

  os.Setenv("CREDHUB_SECRET", cfg.ClientSecret)
  os.Setenv("CREDHUB_CLIENT", cfg.ClientName)
  os.Setenv("CREDHUB_SERVER", cfg.ApiUrl)
  os.Setenv("CREDHUB_CA_CERT", string(uaa_ca_contents)+string(credhub_ca_contents))
}

func TargetAndLoginSkipTls(cfg Config) {
  CleanEnv()

  session := RunCommand("login", "-s", cfg.ApiUrl, "--client-name", cfg.ClientName, "--client-secret", cfg.ClientSecret, "--skip-tls-validation")
  Eventually(session).Should(Exit(0))
}

func CleanEnv() {
  os.Unsetenv("CREDHUB_SECRET")
  os.Unsetenv("CREDHUB_CLIENT")
  os.Unsetenv("CREDHUB_SERVER")
  os.Unsetenv("CREDHUB_CA_CERT")
  os.Unsetenv("CREDHUB_PROXY")
}

type CertificateValue struct {
  Ca          string `yaml:"ca,omitempty"`
  Certificate string `yaml:"certificate,omitempty"`
  PrivateKey  string `yaml:"private_key,omitempty"`
}

type Certificate struct {
  Value CertificateValue `yaml:"value"`
  Name  string           `yaml:"name"`
}
