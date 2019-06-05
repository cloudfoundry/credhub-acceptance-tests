package remote_backend_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
	"strings"
)

var _ = FDescribe("Generate", func() {
	Describe("generate", func() {
		Describe("password type", func() {
			It("generates password type", func() {
				name := "/some-password"

				session := RunCommand("generate", "-t", "password", "-n", name, "-l", "10")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))
				Expect(stdOut).To(ContainSubstring("value: <redacted>"))

				session = RunCommand("get", "-n", name, "-q")
				stdOut = string(session.Out.Contents())
				Expect(strings.TrimSpace(stdOut)).To(HaveLen(10))
			})
		})
		Describe("certificate type", func() {
			It("generates certificate type", func() {
				name := "/some-cert"

				session := RunCommand("generate", "-t", "certificate", "-n", name, "-k", "4096", "-c", "test-cert", "--self-sign")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))
				Expect(stdOut).To(ContainSubstring("value: <redacted>"))

				session = RunCommand("get", "-n", name)
				stdOut = string(session.Out.Contents())
				cert := CertFromPem(stdOut, false)
				Expect(cert.Subject.CommonName).To(Equal("test-cert"))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(4096))
			})
		})
		Describe("rsa type", func() {
			It("generates rsa type", func() {
				name := "/some-rsa"

				session := RunCommand("generate", "-t", "rsa", "-n", name, "-k", "4096")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))
				Expect(stdOut).To(ContainSubstring("value: <redacted>"))

				session = RunCommand("get", "-n", name)
				stdOut = string(session.Out.Contents())
				Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
			})
		})
		Describe("ssh type", func() {
			It("generates ssh type", func() {
				name := "/some-ssh"

				session := RunCommand("generate", "-t", "ssh", "-n", name, "-k", "4096")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))
				Expect(stdOut).To(ContainSubstring("value: <redacted>"))

				session = RunCommand("get", "-n", name)
				stdOut = string(session.Out.Contents())
				Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN SSH PRIVATE KEY-----`))
			})
		})
		Describe("user type", func() {
			It("generates user type", func() {
				name := "/some-user"

				session := RunCommand("generate", "-t", "user", "-n", name, "-l", "47", "-U")
				Expect(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(ContainSubstring(name))
				Expect(stdOut).To(ContainSubstring("value: <redacted>"))

				session = RunCommand("get", "-n", name, "-k", "password", "-q")
				stdOut = string(session.Out.Contents())
				Expect(strings.TrimSpace(stdOut)).To(HaveLen(47))
				Expect(strings.TrimSpace(stdOut)).To(MatchRegexp(`^[a-z0-9_\-]+$`))

			})
		})
	})
})

func CertFromPem(input string, ca bool) *x509.Certificate {
	type certificateValue struct {
		Ca          string `yaml:"ca,omitempty"`
		Certificate string `yaml:"certificate,omitempty"`
	}
	type certificate struct {
		Value certificateValue `yaml:"value"`
	}

	cert := certificate{}
	err := yaml.Unmarshal([]byte(input), &cert)

	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	var pemCert string
	if ca {
		pemCert = cert.Value.Ca
	} else {
		pemCert = cert.Value.Certificate
	}

	block, _ := pem.Decode([]byte(pemCert))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	parsed_cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	return parsed_cert
}
