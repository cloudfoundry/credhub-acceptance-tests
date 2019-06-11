package acceptance_test

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/url"
)

var _ = Describe("Certificates", func() {
	Describe("getting certificate metadata", func() {
		It("gets certificate metadata", func() {
			name := testCredentialPath("some-certificate")

			generateCert := generate.Certificate{
				CommonName: "example.com",
				SelfSign:   true,
			}
			certificate, err := credhubClient.GenerateCertificate(name, generateCert, credhub.Overwrite)
			Expect(err).ToNot(HaveOccurred())
			Expect(certificate.Value.Certificate).ToNot(BeEmpty())
			Expect(certificate.Value.PrivateKey).ToNot(BeEmpty())

			queryParams := url.Values{}
			queryParams.Add("name", certificate.Name)
			data, err := credhubClient.Request(http.MethodGet, "/api/v1/certificates/", queryParams, nil, true)
			Expect(err).ToNot(HaveOccurred())

			dec := json.NewDecoder(data.Body)
			response := make(map[string][]credentials.CertificateMetadata)

			err = dec.Decode(&response)
			Expect(err).ToNot(HaveOccurred())

			metadataArray, _ := response["certificates"]
			expected := metadataArray[0]

			metadata, err := credhubClient.GetAllCertificatesMetadata()
			Expect(err).ToNot(HaveOccurred())

			Expect(metadata).To(ContainElement(expected))

		})
	})
})
