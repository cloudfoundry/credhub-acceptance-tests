package api_integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers/certs"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	config         Config
	credhubCA      []byte
	uaaCA          []byte
	clientCACert   []byte
	clientCAKey    []byte
	certsDir       string
	credentialName string
	appGuid        string
)
var _ = BeforeSuite(func() {
	var err error
	config, err = LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	credhubCA, err = ioutil.ReadFile(filepath.Join(config.CredentialRoot, "server_ca_cert.pem"))
	Expect(err).NotTo(HaveOccurred())

	uaaCA, err = ioutil.ReadFile(filepath.Join(config.UAACa))
	Expect(err).NotTo(HaveOccurred())

	clientCACert, err = ioutil.ReadFile(filepath.Join(config.CredentialRoot, "client_ca_cert.pem"))
	Expect(err).NotTo(HaveOccurred())
	clientCAKey, err = ioutil.ReadFile(filepath.Join(config.CredentialRoot, "client_ca_private.pem"))
	Expect(err).NotTo(HaveOccurred())

	os.Unsetenv("CREDHUB_DEBUG")
})
var _ = Describe("library with mtls authentication", func() {
	const CredhubClientCommonName = "credhub_test_client"

	BeforeEach(func() {
		credentialName = fmt.Sprintf("api-client-mtls-test-%d", time.Now().UnixNano())
		appGuid = uuid.NewString()

		var err error
		certsDir, err = ioutil.TempDir("", "credhub-acceptance-mtls")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(certsDir)).To(Succeed())
	})

	Describe("with a certificate signed by a trusted CA", func() {
		var permissionUuid string
		var adminCredHubClient *credhub.CredHub

		BeforeEach(func() {
			permissions := map[string]interface{}{
				"actor":      "mtls-app:" + appGuid,
				"path":       "/*",
				"operations": []string{"read", "write", "delete"},
			}

			var err error
			adminCredHubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhubCA), string(uaaCA)),
				credhub.Auth(
					auth.UaaClientCredentials(config.ClientName, config.ClientSecret),
				))
			Expect(err).ToNot(HaveOccurred())

			resp, err := adminCredHubClient.Request("POST", "/api/v2/permissions", nil, permissions, false)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			var resBody []byte
			resBody, err = ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			type body struct {
				Uuid string `json:"uuid"`
			}

			var postBody body
			err = json.Unmarshal(resBody, &postBody)
			Expect(err).ToNot(HaveOccurred())
			permissionUuid = postBody.Uuid
			Expect(resp.StatusCode).To(Equal(201))
		})

		AfterEach(func() {
			resp, err := adminCredHubClient.Request("DELETE", "/api/v2/permissions/"+permissionUuid, nil, nil, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))
		})

		It("can do authenticated operations", func() {
			cert, key, err := GenerateSigned(CertOptions{
				CommonName:         CredhubClientCommonName,
				OrganizationalUnit: "app:" + appGuid,
			}, clientCACert, clientCAKey)
			Expect(err).NotTo(HaveOccurred())

			certPath := filepath.Join(certsDir, "cert.pem")
			Expect(ioutil.WriteFile(certPath, cert, 0644)).To(Succeed())
			keyPath := filepath.Join(certsDir, "key.pem")
			Expect(ioutil.WriteFile(keyPath, key, 0600)).To(Succeed())

			credhubClient, err := credhub.New(
				config.ApiUrl,
				credhub.CaCerts(string(credhubCA), string(uaaCA)),
				credhub.ClientCert(certPath, keyPath),
			)
			Expect(err).NotTo(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err = credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("with an expired certificate", func() {
		It("fails on access to authenticated operation ", func() {
			cert, key, err := GenerateSigned(CertOptions{
				CommonName: CredhubClientCommonName,
				NotBefore:  time.Now().Add(time.Hour * 24 * -10),
				NotAfter:   time.Now().Add(time.Hour * 24 * -5),
			}, clientCACert, clientCAKey)
			Expect(err).NotTo(HaveOccurred())

			certPath := filepath.Join(certsDir, "cert.pem")
			Expect(ioutil.WriteFile(certPath, cert, 0644)).To(Succeed())
			keyPath := filepath.Join(certsDir, "key.pem")
			Expect(ioutil.WriteFile(keyPath, key, 0600)).To(Succeed())

			credhubClient, err := credhub.New(
				config.ApiUrl,
				credhub.CaCerts(string(credhubCA), string(uaaCA)),
				credhub.ClientCert(certPath, keyPath),
			)
			Expect(err).NotTo(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err = credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown certificate"))
		})
	})

	Describe("with a self-signed certificate", func() {
		It("fails on access to authenticated operation ", func() {
			cert, key, err := GenerateSelfSigned(CertOptions{})
			Expect(err).NotTo(HaveOccurred())

			certPath := filepath.Join(certsDir, "cert.pem")
			Expect(ioutil.WriteFile(certPath, cert, 0644)).To(Succeed())
			keyPath := filepath.Join(certsDir, "key.pem")
			Expect(ioutil.WriteFile(keyPath, key, 0600)).To(Succeed())

			credhubClient, err := credhub.New(
				config.ApiUrl,
				credhub.CaCerts(string(credhubCA), string(uaaCA)),
				credhub.ClientCert(certPath, keyPath),
			)
			Expect(err).NotTo(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err = credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)

			Expect(err.Error()).To(Equal("invalid_token: Full authentication is required to access this resource"))
		})
	})

	Describe("with certificate signed by unknown CA", func() {
		It("fails on access to authenticated operation ", func() {
			caCert, caKey, err := GenerateSelfSigned(CertOptions{})
			Expect(err).NotTo(HaveOccurred())
			cert, key, err := GenerateSigned(CertOptions{}, caCert, caKey)
			Expect(err).NotTo(HaveOccurred())

			certPath := filepath.Join(certsDir, "cert.pem")
			Expect(ioutil.WriteFile(certPath, cert, 0644)).To(Succeed())
			keyPath := filepath.Join(certsDir, "key.pem")
			Expect(ioutil.WriteFile(keyPath, key, 0600)).To(Succeed())

			credhubClient, err := credhub.New(
				config.ApiUrl,
				credhub.CaCerts(string(credhubCA), string(uaaCA)),
				credhub.ClientCert(certPath, keyPath),
			)
			Expect(err).NotTo(HaveOccurred())

			generatePassword := generate.Password{Length: 10}
			_, err = credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)

			Expect(err.Error()).To(Equal("invalid_token: Full authentication is required to access this resource"))
		})
	})
})

func TestLibraryMTLS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mTLS API Library Test Suite")
}
