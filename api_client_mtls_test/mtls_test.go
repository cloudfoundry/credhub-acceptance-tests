package api_integration_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var (
	credhubClient  *credhub.CredHub
	config         Config
	err            error
	credhub_ca     []byte
	uaa_ca         []byte
	certPath       string
	credentialName string
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
		credentialName = fmt.Sprintf("%d", time.Now().UnixNano())
	})

	Describe("with a certificate signed by a trusted CA", func() {
		var permissionUuid string
		var adminCredHubClient *credhub.CredHub

		BeforeEach(func() {
			certBytes, err := ioutil.ReadFile(path.Join(certPath, "client.pem"))
			Expect(err).ToNot(HaveOccurred())

			block, pemByte := pem.Decode([]byte(strings.TrimSpace(string(certBytes))))
			pemCertsByte := append(block.Bytes, pemByte...)
			cert, err := x509.ParseCertificate(pemCertsByte)
			Expect(err).ToNot(HaveOccurred())

			userID := "mtls-" + cert.Subject.OrganizationalUnit[0]
			vals := map[string]interface{}{
				"actor":      userID,
				"path":       "/*",
				"operations": []string{"read", "write", "delete"},
			}

			adminCredHubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
				credhub.Auth(
					auth.UaaClientCredentials(config.ClientName, config.ClientSecret),
				))
			Expect(err).ToNot(HaveOccurred())

			resp, err := adminCredHubClient.Request("POST", "/api/v2/permissions", nil, vals, false)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			var resBody []byte
			resBody, err = ioutil.ReadAll(resp.Body)

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
			credhubClient, err = credhub.New(config.ApiUrl,
				credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
				credhub.ClientCert(path.Join(certPath, "client.pem"),
					path.Join(certPath, "client_key.pem")))

			Expect(err).ToNot(HaveOccurred())
			generatePassword := generate.Password{Length: 10}
			_, err = credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)
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
			_, err := credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)
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
			_, err := credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)

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
			_, err := credhubClient.GeneratePassword(credentialName, generatePassword, credhub.Overwrite)

			Expect(err.Error()).To(Equal("invalid_token: Full authentication is required to access this resource"))
		})

	})
})

func TestLibraryMTLS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mTLS API Library Test Suite")
}
