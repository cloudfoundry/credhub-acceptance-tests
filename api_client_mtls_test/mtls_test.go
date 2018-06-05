package api_integration_test

import (
	"io/ioutil"
	"path"
	"testing"

	"os"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	credhubClient *credhub.CredHub
	config        Config
	err           error
	credhub_ca    []byte
	uaa_ca        []byte
	certPath      string
)

var _ = Describe("library with mtls authentication", func() {

	BeforeEach(func() {
		config, err = LoadConfig()
		Expect(err).NotTo(HaveOccurred())
		credhub_ca, err = ioutil.ReadFile(path.Join(config.CredentialRoot, "server_ca_cert.pem"))
		Expect(err).NotTo(HaveOccurred())
		uaa_ca, err = ioutil.ReadFile(path.Join(config.UAACa))
		Expect(err).NotTo(HaveOccurred())
		certPath = path.Join(os.Getenv("PWD"), "certs")
	})

	Describe("with a certificate signed by a trusted CA", func() {
		It("can do authenticated operations", func() {
			credhubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
				credhub.ClientCert(path.Join(certPath, "client.pem"),
					path.Join(certPath, "client_key.pem")))

			Expect(err).ToNot(HaveOccurred())
			generatePassword := generate.Password{Length: 10}
			_, err := credhubClient.GeneratePassword("test", generatePassword, credhub.Overwrite)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("with an expired certificate", func() {
		It("fails on access to authenticated operation ", func() {
			credhubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
				credhub.ClientCert(path.Join(certPath, "expired.pem"),
					path.Join(certPath, "expired_key.pem")))

			Expect(err).ToNot(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err := credhubClient.GeneratePassword("test", generatePassword, credhub.Overwrite)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown certificate"))
		})
	})

	Describe("with a self-signed certificate", func() {
		It("fails on access to authenticated operation ", func() {
			credhubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
				credhub.ClientCert(path.Join(certPath, "selfsigned.pem"),
					path.Join(certPath, "selfsigned_key.pem")))

			Expect(err).ToNot(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err := credhubClient.GeneratePassword("test", generatePassword, credhub.Overwrite)

			Expect(err.Error()).To(Equal("invalid_token: Full authentication is required to access this resource"))
		})

	})

	Describe("with certificate signed by unknown CA", func() {
		It("fails on access to authenticated operation ", func() {
			credhubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
				credhub.ClientCert(path.Join(certPath, "unknown.pem"),
					path.Join(certPath, "unknown_key.pem")))

			Expect(err).ToNot(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err := credhubClient.GeneratePassword("test", generatePassword, credhub.Overwrite)

			Expect(err.Error()).To(Equal("invalid_token: Full authentication is required to access this resource"))
		})

	})
})

func TestLibraryMTLS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mTLS API Library Test Suite")
}
