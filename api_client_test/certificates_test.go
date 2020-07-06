package acceptance_test

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/url"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"fmt"
	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Certificates", func() {
	Describe("getting certificate metadata", func() {
		It("gets certificate metadata", func() {
			name := testCredentialPath(time.Now().UnixNano(), "some-certificate")

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

			actual, err := credhubClient.GetAllCertificatesMetadata()
			Expect(err).ToNot(HaveOccurred())

			Expect(actual).To(ContainElement(expected))

		})
		It("properly returns self_signed and is_ca", func() {
			name := testCredentialPath(time.Now().UnixNano(), "some-intermediate-ca")

			setCertificate := values.Certificate{
				Certificate: test_helpers.INTERMEDIATE_CA,
				PrivateKey:  test_helpers.INTERMEDIATE_CA_PRIVATE_KEY,
			}
			_, err := credhubClient.SetCertificate(name, setCertificate)
			Expect(err).ToNot(HaveOccurred())

			queryParams := url.Values{}
			queryParams.Add("name", name)
			data, err := credhubClient.Request(http.MethodGet, "/api/v1/certificates/", queryParams, nil, true)
			Expect(err).ToNot(HaveOccurred())

			dec := json.NewDecoder(data.Body)
			response := make(map[string][]credentials.CertificateMetadata)

			err = dec.Decode(&response)
			Expect(err).ToNot(HaveOccurred())

			metadataArray, _ := response["certificates"]
			actual := metadataArray[0]

			Expect(actual.Versions[0].SelfSigned).To(BeFalse())
			Expect(actual.Versions[0].CertificateAuthority).To(BeTrue())
		})
	})
	Describe("update transitional version to latest", func() {
		It("it accepts latest in the request body", func() {
			name := testCredentialPath(time.Now().UnixNano(), "some-certificate")

			generateCert := generate.Certificate{
				CommonName: "example.com",
				SelfSign:   true,
				IsCA:       true,
			}
			certificate, err := credhubClient.GenerateCertificate(name, generateCert, credhub.Overwrite)
			Expect(err).ToNot(HaveOccurred())
			certificate2, err := credhubClient.Regenerate(name)
			Expect(err).ToNot(HaveOccurred())

			queryParams := url.Values{}
			queryParams.Add("name", name)
			data, err := credhubClient.Request(http.MethodGet, "/api/v1/certificates/", queryParams, nil, true)
			Expect(err).ToNot(HaveOccurred())

			dec := json.NewDecoder(data.Body)
			response := make(map[string][]credentials.CertificateMetadata)

			err = dec.Decode(&response)
			Expect(err).ToNot(HaveOccurred())

			metadataArray, _ := response["certificates"]
			cert := metadataArray[0]

			requestBody := map[string]interface{}{
				"version": "latest",
			}
			pathString := fmt.Sprintf("/api/v1/certificates/%s/update_transitional_version", cert.Id)
			_, err = credhubClient.Request(http.MethodPut, pathString, nil, requestBody, true)
			Expect(err).ToNot(HaveOccurred())

			data, err = credhubClient.Request(http.MethodGet, "/api/v1/certificates/", queryParams, nil, true)
			Expect(err).ToNot(HaveOccurred())

			dec = json.NewDecoder(data.Body)
			response = make(map[string][]credentials.CertificateMetadata)

			err = dec.Decode(&response)
			Expect(err).ToNot(HaveOccurred())

			metadataArray, _ = response["certificates"]
			actual := metadataArray[0]

			Expect(len(actual.Versions)).To(Equal(2))
			Expect(actual.Versions[0].Id).To(Equal(certificate2.Id))
			Expect(actual.Versions[1].Id).To(Equal(certificate.Id))

			Expect(actual.Versions[0].Transitional).To(BeTrue())
			Expect(actual.Versions[1].Transitional).To(BeFalse())

		})
	})

	Describe("regenerate allow_transitional_parent_to_sign is present", func() {
		When("transitional version is latest", func() {
			It("can sign certificates with a transitional version", func() {
				//setup ca generation
				caName := testCredentialPath(time.Now().UnixNano(), "some-ca")
				caGenerationParams := generate.Certificate{
					CommonName: "some-ca",
					SelfSign:   true,
					IsCA:       true,
				}

				//generate a ca
				_, err := credhubClient.GenerateCertificate(caName, caGenerationParams, credhub.Overwrite)
				Expect(err).ToNot(HaveOccurred())

				//setup leaf generation
				certName := testCredentialPath(time.Now().UnixNano(), "some-cert")
				certGenerationParams := generate.Certificate{
					CommonName: "some-cert",
					Ca:   caName,
				}

				//generate a leaf
				_, err = credhubClient.GenerateCertificate(certName, certGenerationParams, credhub.Overwrite)
				Expect(err).ToNot(HaveOccurred())

				//get certificate ids
				data, err := credhubClient.Request(http.MethodGet, "/api/v1/certificates/", nil, nil, true)
				Expect(err).ToNot(HaveOccurred())

				dec := json.NewDecoder(data.Body)
				response := make(map[string][]credentials.CertificateMetadata)

				err = dec.Decode(&response)
				Expect(err).ToNot(HaveOccurred())

				ids := map[string]string{}
				metadataArray, _ := response["certificates"]
				for _, item := range metadataArray {
					ids[item.Name] = item.Id
				}

				//regenerate ca, setting as transitional
				requestBody := map[string]interface{}{
					"set_as_transitional": true,
				}
				pathString := fmt.Sprintf("/api/v1/certificates/%s/regenerate", ids[caName])
				r, err := credhubClient.Request(http.MethodPost, pathString, nil, requestBody, true)
				Expect(err).ToNot(HaveOccurred())

				var regeneratedCa credentials.Certificate
				dec = json.NewDecoder(r.Body)

				err = dec.Decode(&regeneratedCa)
				Expect(err).ToNot(HaveOccurred())


				regeneratedCaPem := fmt.Sprintf("%v", regeneratedCa.Value.Certificate)

				//regenerate leaf, using allow_transitional_parent_to_sign
				requestBody = map[string]interface{}{
					"allow_transitional_parent_to_sign": true,
				}
				pathString = fmt.Sprintf("/api/v1/certificates/%s/regenerate", ids[certName])
				r, err = credhubClient.Request(http.MethodPost, pathString, nil, requestBody, true)
				Expect(err).ToNot(HaveOccurred())

				var regeneratedCert credentials.Certificate
				dec = json.NewDecoder(r.Body)

				err = dec.Decode(&regeneratedCert)
				Expect(err).ToNot(HaveOccurred())
				regeneratedCertPem := fmt.Sprintf("%v", regeneratedCert.Value.Certificate)


				//check if transitional ca version signs cert
				decodedCa, _ := pem.Decode([]byte(regeneratedCaPem))
				parsedCa, err := x509.ParseCertificate(decodedCa.Bytes)
				Expect(err).ToNot(HaveOccurred())
				pool := x509.NewCertPool()
				pool.AddCert(parsedCa)

				decodedCert, _ := pem.Decode([]byte(regeneratedCertPem))
				parsedCert, err := x509.ParseCertificate(decodedCert.Bytes)
				Expect(err).ToNot(HaveOccurred())
				_, err = parsedCert.Verify(x509.VerifyOptions{Roots: pool})
				Expect(err).ToNot(HaveOccurred())

			})
		})
	})
})
