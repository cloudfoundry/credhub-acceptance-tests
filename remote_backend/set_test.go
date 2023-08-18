package remote_backend_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Set", func() {
	Describe("set", func() {
		Describe("json type", func() {
			It("sets json type", func() {
				name := GenerateUniqueCredentialName()
				json := `{"some-key": "some-json"}`
				asYaml := "some-key: some-json"

				session := RunCommand("set", "-t", "json", "-n", name, "-v", json)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(asYaml))
				Expect(stdOut).To(ContainSubstring(name))
			})
		})
		Describe("value type", func() {
			It("sets value type", func() {
				name := GenerateUniqueCredentialName()
				value := "some-random-value"

				session := RunCommand("set", "-t", "value", "-n", name, "-v", value)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(value))
				Expect(stdOut).To(ContainSubstring(name))
			})
		})
		Describe("password type", func() {
			It("sets password type", func() {
				name := GenerateUniqueCredentialName()
				password := "some-super-secret-password"

				session := RunCommand("set", "-t", "password", "-n", name, "-w", password)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(password))
				Expect(stdOut).To(ContainSubstring(name))
			})
		})
		Describe("certificate type", func() {
			It("sets certificate type", func() {
				name := GenerateUniqueCredentialName()
				cert := VALID_CERTIFICATE
				privateKey := VALID_CERTIFICATE_PRIVATE_KEY
				ca := VALID_CERTIFICATE_CA

				session := RunCommand("set", "-t", "certificate",
					"-n", name,
					"-r", ca,
					"-c", cert,
					"-p", privateKey)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))

				type getCertificateResponse struct {
					Value struct {
						CA          string `yaml:"ca"`
						Certificate string `yaml:"certificate"`
						PrivateKey  string `yaml:"private_key"`
					} `yaml:"value"`
				}
				var rsp getCertificateResponse
				Expect(yaml.Unmarshal([]byte(stdOut), &rsp)).To(Succeed())
				Expect(rsp.Value.CA).To(Equal(ca))
				Expect(rsp.Value.Certificate).To(Equal(cert))
				Expect(rsp.Value.PrivateKey).To(Equal(privateKey))
			})
		})
		Describe("user type", func() {
			It("sets user type", func() {
				name := GenerateUniqueCredentialName()
				username := "some-username"
				password := "some-random-password"

				session := RunCommand("set", "-t", "user", "-n", name, "-z", username, "-w", password)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(username))
				Expect(stdOut).To(ContainSubstring(password))
				Expect(stdOut).To(ContainSubstring(name))
			})
		})
		Describe("rsa type", func() {
			It("sets rsa type", func() {
				name := GenerateUniqueCredentialName()
				publicKey := ALTERNATE_CA_PUBLIC_KEY
				privateKey := ALTERNATE_CA_PRIVATE_KEY

				session := RunCommand("set", "-t", "rsa", "-n", name, "-u", publicKey, "-p", privateKey)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))

				type getRSAResponse struct {
					Value struct {
						PublicKey  string `yaml:"public_key"`
						PrivateKey string `yaml:"private_key"`
					} `yaml:"value"`
				}
				var rsp getRSAResponse
				Expect(yaml.Unmarshal([]byte(stdOut), &rsp)).To(Succeed())
				Expect(rsp.Value.PublicKey).To(Equal(publicKey))
				Expect(rsp.Value.PrivateKey).To(Equal(privateKey))
			})
		})
		Describe("ssh type", func() {
			It("sets ssh type", func() {
				name := GenerateUniqueCredentialName()
				publicKey := ALTERNATE_CA_PUBLIC_KEY
				privateKey := ALTERNATE_CA_PRIVATE_KEY

				session := RunCommand("set", "-t", "ssh", "-n", name, "-u", publicKey, "-p", privateKey)
				Expect(session).Should(Exit(0))

				session = RunCommand("get", "-n", name)
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))

				type getSSHResponse struct {
					Value struct {
						PublicKey  string `yaml:"public_key"`
						PrivateKey string `yaml:"private_key"`
					} `yaml:"value"`
				}
				var rsp getSSHResponse
				Expect(yaml.Unmarshal([]byte(stdOut), &rsp)).To(Succeed())
				Expect(rsp.Value.PublicKey).To(Equal(publicKey))
				Expect(rsp.Value.PrivateKey).To(Equal(privateKey))
			})
		})
	})
})
