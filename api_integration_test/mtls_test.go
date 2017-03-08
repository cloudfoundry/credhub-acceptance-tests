package api_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"fmt"
	"crypto/tls"
	"log"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
	"testing"
)

var _ = Describe("mutual TLS authentication", func() {
	var (
		config Config
		err error
	)

	BeforeEach(func() {
		config, err = LoadConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("with a certificate signed by a trusted CA	", func() {
		Describe("when the certificate has a valid date range", func() {
			It("allows the user to hit an authenticated endpoint", func() {
				result := parlezMTLS(
						config.ApiUrl + "/api/v1/data",
						config.ValidClientCertPath,
						config.ValidClientKeyPath,
						config.ValidServerCAPath)

				Expect(result).To(MatchRegexp(`"type":"password"`))
			})
		})
	})
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

func parlezMTLS(url, clientCertPath, clientKeyPath, serverCAPath string) string {
	Expect(clientCertPath).NotTo(BeEmpty())
	Expect(clientKeyPath).NotTo(BeEmpty())
	Expect(serverCAPath).NotTo(BeEmpty())

	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	handleError(err)

	roots := x509.NewCertPool()

	CA, err := ioutil.ReadFile(serverCAPath)
	ok := roots.AppendCertsFromPEM([]byte(CA))
	if !ok {
		panic("failed to parse root certificate")
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: roots,
	}
	tr := &http.Transport{TLSClientConfig: tlsConf}
	client := &http.Client{Transport: tr}

	values := map[string]string{"name": "mtlstest", "type": "password"}
	jsonValue, _ := json.Marshal(values)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	handleError(err)

	fmt.Println(resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	return string(body)
}
