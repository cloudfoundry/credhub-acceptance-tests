package api_integration_test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers/certs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("mutual TLS authentication", func() {
	const CredhubClientCommonName = "credhub_test_client"

	var (
		config         Config
		credhubCA      []byte
		uaaCA          []byte
		clientCACert   []byte
		clientCAKey    []byte
		credentialName string
		appGuid        string
	)

	BeforeSuite(func() {
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
	})

	BeforeEach(func() {
		credentialName = fmt.Sprintf("api-integration-test-%d", time.Now().UnixNano())
		appGuid = uuid.NewV4().String()
	})

	Describe("with a certificate signed by a trusted CA", func() {
		var (
			adminCredHubClient *credhub.CredHub
			permissionUuid     string
		)

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
			Expect(err).NotTo(HaveOccurred())

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

		It("allows the client to hit an authenticated endpoint", func() {
			cert, key, err := GenerateSigned(CertOptions{
				CommonName:         CredhubClientCommonName,
				OrganizationalUnit: "app:" + appGuid,
			}, clientCACert, clientCAKey)
			Expect(err).NotTo(HaveOccurred())

			postData := map[string]string{"name": credentialName, "type": "password"}
			result, err := mtlsPost(config.ApiUrl+"/api/v1/data", postData, credhubCA, cert, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(MatchRegexp(`"type":"password"`))
		})
	})

	Describe("with an expired certificate", func() {
		It("prevents the client from hitting an authenticated endpoint", func() {
			cert, key, err := GenerateSigned(CertOptions{
				CommonName: CredhubClientCommonName,
				NotBefore:  time.Now().Add(time.Hour * 24 * -10),
				NotAfter:   time.Now().Add(time.Hour * 24 * -5),
			}, clientCACert, clientCAKey)
			Expect(err).NotTo(HaveOccurred())

			postData := map[string]string{"name": credentialName, "type": "password"}
			result, err := mtlsPost(config.ApiUrl+"/api/v1/data", postData, credhubCA, cert, key)
			Expect(err).To(MatchError(ContainSubstring("unknown certificate")))
			Expect(result).To(BeEmpty())
		})
	})

	Describe("with a self-signed certificate", func() {
		It("prevents the client from hitting an authenticated endpoint", func() {
			cert, key, err := GenerateSelfSigned(CertOptions{})
			Expect(err).NotTo(HaveOccurred())

			postData := map[string]string{"name": credentialName, "type": "password"}
			result, err := mtlsPost(config.ApiUrl+"/api/v1/data", postData, credhubCA, cert, key)

			// golang doesn't seem to send self-signed certs
			// server.ssl.client-auth=want (https://tools.ietf.org/html/rfc5246#section-7.4.4)
			// That is why, we are asserting on OAuth authorization failure here.
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(MatchRegexp(".*Full authentication is required to access this resource"))
		})
	})

	Describe("with a certificate signed by an unknown CA", func() {
		It("prevents the client from hitting an authenticated endpoint", func() {
			caCert, caKey, err := GenerateSelfSigned(CertOptions{})
			Expect(err).NotTo(HaveOccurred())
			cert, key, err := GenerateSigned(CertOptions{}, caCert, caKey)
			Expect(err).NotTo(HaveOccurred())

			postData := map[string]string{"name": credentialName, "type": "password"}
			result, err := mtlsPost(config.ApiUrl+"/api/v1/data", postData, credhubCA, cert, key)

			// Okay, so golang 1.7.x **sometimes** doesn't seem to send certs that the server won't accept...
			// Here we assert that, if there was an error, it should be the server rejecting the cert, and
			// if there wasn't an error, the server told us to go away because we didn't send an auth token or cert.
			if err != nil {
				Expect(err.Error()).To(ContainSubstring("unknown certificate"))
			} else {
				Expect(result).To(MatchRegexp(".*Full authentication is required to access this resource"))
			}
		})
	})
})

func TestMTLS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mTLS Test Suite")
}

func mtlsPost(url string, postData map[string]string, serverCA, clientCert, clientKey []byte) (string, error) {
	clientCertificate, err := tls.X509KeyPair(clientCert, clientKey)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	trustedCAs := x509.NewCertPool()
	ok := trustedCAs.AppendCertsFromPEM([]byte(serverCA))
	if !ok {
		return "", errors.New("failed to parse server CA")
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{clientCertificate},
		RootCAs:      trustedCAs,
	}

	transport := &http.Transport{TLSClientConfig: tlsConf}
	client := &http.Client{Transport: transport}

	jsonValue, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
