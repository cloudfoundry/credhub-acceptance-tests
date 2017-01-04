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

	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	commandPath string
	homeDir     string
	cfg         Config
	err         error
)

// We look for these values in the verify-logging CI task to ensure that credentials don't leak
const credentialValue = "FAKE-CREDENTIAL-VALUE"
const newCredentialValue = "FAKE-CREDENTIAL-VALUE1"

var _ = Describe("Integration test", func() {
	BeforeEach(func() {
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

	It("should set, get, and delete a new value secret", func() {
		credentialName := generateUniqueCredentialName()

		By("trying to access a secret that doesn't exist", func() {
			session := runCommand("get", "-n", credentialName)
			stdErr := string(session.Err.Contents())

			Expect(stdErr).To(MatchRegexp(`Credential not found. Please validate your input and retry your request.`))
			Eventually(session).Should(Exit(1))
		})

		By("setting a new value secret", func() {
			session := runCommand("set", "-n", credentialName, "-t", "value", "-v", credentialValue)
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
		})

		By("getting the new value secret", func() {
			session := runCommand("get", "-n", credentialName)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
		})

		By("deleting the secret", func() {
			session := runCommand("delete", "-n", credentialName)
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("handling special characters", func() {
		It("should handle secrets who names begin with a leading slash", func() {
			baseId := "ace/ventura" + generateUniqueCredentialName()
			leadingSlashId := "/" + baseId
			passwordValue := "finkel-is-einhorn"

			By("setting a value whose name begins with a leading slash", func() {
				session := runCommand("set", "-n", leadingSlashId, "-t", "password", "-v", passwordValue)
				Eventually(session).Should(Exit(0))
			})

			By("retrieving the value that was set with a leading slash", func() {
				session := runCommand("get", "-n", leadingSlashId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
				Expect(stdOut).To(ContainSubstring(passwordValue))
			})

			By("retrieving the value that was set without a leading slash", func() {
				session := runCommand("get", "-n", baseId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
				Expect(stdOut).To(ContainSubstring(passwordValue))
			})
		})

		It("should get secrets whose names have lots of special characters", func() {
			crazyCharsId := "dan:test/ing?danother[stuff]that@shouldn!tbe$in&the" + generateUniqueCredentialName()

			By("setting a value with lots of special characters", func() {
				session := runCommand("set", "-n", crazyCharsId, "-t", "password", "-v", "woof-woof")
				Eventually(session).Should(Exit(0))
			})

			By("retrieving the value that was set", func() {
				session := runCommand("get", "-n", crazyCharsId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
				Expect(stdOut).To(ContainSubstring(crazyCharsId))
			})
		})

		It("should handle edge-casey character combinations", func() {
			edgeCaseId := "&gunk=x/bar/cr@zytown108" + generateUniqueCredentialName()

			By("setting a value with lots of special characters", func() {
				session := runCommand("set", "-n", edgeCaseId, "-t", "password", "-v", "find-me")
				Eventually(session).Should(Exit(0))
			})

			By("retrieving the value that was set", func() {
				session := runCommand("get", "-n", edgeCaseId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
				Expect(stdOut).To(ContainSubstring(edgeCaseId))
			})
		})

		It("should delete secrets with special characters", func() {
			deleteId := "?testParam=foo&gunk=x/bar/piv0t@l" + generateUniqueCredentialName()

			By("setting a value with lots of special characters", func() {
				session := runCommand("set", "-n", deleteId, "-t", "password", "-v", "find-me")
				Eventually(session).Should(Exit(0))
			})

			By("deleting the secret", func() {
				session := runCommand("delete", "-n", deleteId)
				Eventually(session).Should(Exit(0))
			})
		})
	})

	It("should generate a password", func() {
		session := runCommand("generate", "-n", generateUniqueCredentialName(), "-t", "password")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
	})

	It("should regenerate passwords with similar rules", func() {
		generatedPasswordId := generateUniqueCredentialName()

		By("first generating a password with no numbers", func() {
			session := runCommand("generate", "-n", generatedPasswordId, "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(MatchRegexp(`Type:\s+password`))
			Expect(stdOut).To(Not(MatchRegexp(`Value:\s.+\d`)))
		})

		By("then regenerating the password and observing it still has no numbers", func() {
			session := runCommand("regenerate", "-n", generatedPasswordId)
			Eventually(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(Not(MatchRegexp(`Value:\s.+\d`)))
		})
	})

	It("should set a secret's timestamp correctly", func() {
		var original_timestamp []byte
		r, _ := regexp.Compile(`Updated:\s+(.*)[\s|$]`)

		credentialName := generateUniqueCredentialName()

		By("getting the original timestamp", func() {
			session := runCommand("set", "-n", credentialName, "-t", "value", "-v", credentialValue)
			original_timestamp_array := r.FindSubmatch(session.Out.Contents())

			Expect(original_timestamp_array).To(HaveLen(2))

			original_timestamp = original_timestamp_array[1]

			Expect(original_timestamp).NotTo(HaveLen(0))
		})

		By("getting the timestamp after a no-overwrite set", func() {
			session := runCommand("set", "-n", credentialName, "-t", "value", "-v", credentialValue, "--no-overwrite")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp("Value:\\s+" + credentialValue))
			Expect(stdOut).To(MatchRegexp(fmt.Sprintf(`Updated:\s+%s`, original_timestamp)))
		})

		By("getting the timestamp after an overwrite set", func() {
			// We need to sleep in order to ensure that the timestamp is different,
			// since it is truncated to the second.
			time.Sleep(time.Duration(1) * time.Second)

			session := runCommand("set", "-n", credentialName, "-t", "value", "-v", newCredentialValue)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp("Value:\\s+" + newCredentialValue))
			Expect(stdOut).NotTo(MatchRegexp(fmt.Sprintf(`Updated:\s+%s`, original_timestamp)))
		})

		By("getting the value", func() {
			session := runCommand("get", "-n", credentialName)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
			Expect(stdOut).To(MatchRegexp("Value:\\s+" + newCredentialValue))
		})
	})

	Describe("setting a certificate", func() {
		It("should be able to set a certificate", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "certificate", "--certificate-string", "iamacertificate")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
			Expect(stdOut).To(MatchRegexp(`Certificate:\s+iamacertificate`))
		})

		It("should require a certificate type", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "certificate", "--no-overwrite")
			Eventually(session).Should(Exit(1))
			Expect(session.Err.Contents()).To(MatchRegexp(".*At least one certificate type must be set. Please validate your input and retry your request."))
		})
	})

	Describe("setting an SSH key", func() {
		It("should be able to set an ssh key", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "ssh", "-U", "iamapublickey", "-P", credentialValue)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+iamapublickey`))
			Expect(stdOut).To(MatchRegexp("Private Key:\\s+" + credentialValue))
		})
	})

	Describe("setting an RSA key", func() {
		It("should be able to set an rsa key", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "rsa", "-U", "iamapublickey", "-P", credentialValue)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+rsa`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+iamapublickey`))
			Expect(stdOut).To(MatchRegexp("Private Key:\\s+" + credentialValue))
		})
	})

	Describe("CAs and Certificates", func() {
		It("should generate a CA and certificate", func() {
			certificateAuthorityId := generateUniqueCredentialName()
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
				session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_signing")
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

		It("should error gracefully when supplying an invalid extended key name", func() {
			certificateAuthorityId := generateUniqueCredentialName()
			certificateId := certificateAuthorityId + "1"
			runCommand("ca-generate", "-n", certificateAuthorityId, "--common-name", certificateAuthorityId)
			session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_sinning")
			stdErr := string(session.Err.Contents())

			Eventually(session).Should(Exit(1))
			Expect(stdErr).To(MatchRegexp(`The provided extended key usage 'code_sinning' was not known. Please update this value and retry your request.`))
		})

		It("should handle CAs whose names have lots of special characters", func() {
			madDogCAId := "dan:test/ing?danothertbe$in&the[stuff]=that@shouldn!"

			By("setting a value with lots of special characters", func() {
				session := runCommand("ca-generate", "-n", madDogCAId, "--common-name", generateUniqueCredentialName())
				Eventually(session).Should(Exit(0))
			})

			By("retrieving the value that was set", func() {
				session := runCommand("ca-get", "-n", madDogCAId)
				Eventually(session).Should(Exit(0))
			})
		})
	})

	It("should generate an SSH key", func() {
		sshSecretName := generateUniqueCredentialName()

		By("generating the key", func() {
			session := runCommand("generate", "-n", sshSecretName, "-t", "ssh")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+ssh`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+ssh-rsa \S+`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := runCommand("get", "-n", sshSecretName)
			Eventually(session).Should(Exit(0))
		})
	})

	It("should generate an RSA key", func() {
		rsaSecretName := generateUniqueCredentialName()

		By("generating the key", func() {
			session := runCommand("generate", "-n", rsaSecretName, "-t", "rsa")

			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())

			Expect(stdOut).To(MatchRegexp(`Type:\s+rsa`))
			Expect(stdOut).To(MatchRegexp(`Public Key:\s+-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
		})

		By("getting the key", func() {
			session := runCommand("get", "-n", rsaSecretName)
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
