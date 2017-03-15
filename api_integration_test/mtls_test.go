package api_integration_test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"testing"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	config Config
	err    error
)

var _ = Describe("mutual TLS authentication", func() {

	Describe("with a certificate signed by a trusted CA	", func() {
		BeforeEach(func() {
			config, err = LoadConfig()
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows the client to hit an authenticated endpoint", func() {
			postData := map[string]string{
				"name": "mtlstest",
				"type": "password",
			}
			result, err := mtlsPost(
				config.ApiUrl+"/api/v1/data",
				postData,
				"server_ca_cert.pem",
				"client.pem",
				"client_key.pem")

			Expect(err).To(BeNil())
			Expect(result).To(MatchRegexp(`"type":"password"`))
		})
	})

	Describe("with an expired certificate", func() {
		BeforeEach(func() {
			config, err = LoadConfig()
			Expect(err).NotTo(HaveOccurred())
		})

		It("prevents the client from hitting an authenticated endpoint", func() {
			postData := map[string]string{
				"name": "mtlstest",
				"type": "password",
			}
			result, err := mtlsPost(
				config.ApiUrl+"/api/v1/data",
				postData,
				"server_ca_cert.pem",
				"expired.pem",
				"expired_key.pem")

			Expect(err).ToNot(BeNil())
			Expect(result).To(BeEmpty())
		})
	})

	Describe("with a self-signed certificate", func() {
		BeforeEach(func() {
			config, err = LoadConfig()
			Expect(err).NotTo(HaveOccurred())
		})

		It("prevents the client from hitting an authenticated endpoint", func() {
			postData := map[string]string{
				"name": "mtlstest",
				"type": "password",
			}
			result, err := mtlsPost(
				config.ApiUrl+"/api/v1/data",
				postData,
				"server_ca_cert.pem",
				"selfsigned.pem",
				"selfsigned_key.pem")

			// golang doesn't seem to send self-signed certs
			// server.ssl.client-auth=want (https://tools.ietf.org/html/rfc5246#section-7.4.4)
			// That is why, we are asserting on OAuth authorization failure here.
			Expect(err).To(BeNil())
			Expect(result).To(MatchRegexp(".*Full authentication is required to access this resource"))
		})
	})

	//Describe("with a certificate signed by an unknown CA", func() {
	//	BeforeEach(func() {
	//		config, err = LoadConfig()
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//
	//	It("prevents the client from hitting an authenticated endpoint", func() {
	//		postData := map[string]string{
	//			"name": "mtlstest",
	//			"type": "password",
	//		}
	//		result, err := mtlsPost(
	//			config.ApiUrl+"/api/v1/data",
	//			postData,
	//			"server_ca_cert.pem",
	//			"unknown.pem",
	//			"unknown_key.pem")
	//
	//		//Expect(err).ToNot(BeNil())
	//		fmt.Println("The error:")
	//		fmt.Println(err)
	//		fmt.Println(result)
	//		Expect(result).To(BeEmpty())
	//	})
	//})
})

func TestMTLS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mTLS Test Suite")
}

func handleError(err error) {
	if err != nil {
		log.Fatal("Fatal", err)
	}
}

func mtlsPost(url string, postData map[string]string, serverCaFilename, clientCertFilename, clientKeyPath string) (string, error) {
	client, err := createMtlsClient(serverCaFilename, clientCertFilename, clientKeyPath)

	jsonValue, _ := json.Marshal(postData)

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

func createMtlsClient(serverCaFilename, clientCertFilename, clientKeyFilename string) (*http.Client, error) {
	serverCaPath := path.Join(config.CredentialRoot, serverCaFilename)
	clientCertPath := path.Join(os.Getenv("PWD"), "certs", clientCertFilename)
	clientKeyPath := path.Join(os.Getenv("PWD"), "certs", clientKeyFilename)

	_, err := os.Stat(serverCaPath)
	handleError(err)
	_, err = os.Stat(clientCertPath)
	handleError(err)
	_, err = os.Stat(clientKeyPath)
	handleError(err)

	clientCertificate, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	handleError(err)

	trustedCAs := x509.NewCertPool()
	serverCA, err := ioutil.ReadFile(serverCaPath)

	ok := trustedCAs.AppendCertsFromPEM([]byte(serverCA))
	if !ok {
		log.Fatal("failed to parse root certificate")
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{clientCertificate},
		RootCAs:      trustedCAs,
	}

	transport := &http.Transport{TLSClientConfig: tlsConf}
	client := &http.Client{Transport: transport}

	return client, err
}
