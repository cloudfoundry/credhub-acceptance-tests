package integration_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
)

type getCertificatesResponse struct {
	Certificates []struct {
		Id       string               `json:"id"`
		Name     string               `json:"name"`
		Versions []certificateVersion `json:"versions"`
		SignedBy string               `json:"signed_by"`
		Signs    []string             `json:"signs"`
	} `json:"certificates"`
}

type certificateVersion struct {
	Id           string `json:"id"`
	ExpiryDate   string `json:"expiry_date"`
	Transitional bool   `json:"transitional"`
	Generated    bool   `json:"generated"`
}

type getCertificateVersionsResponse []struct {
	Id               string                  `json:"id"`
	Name             string                  `json:"name"`
	ExpiryDate       string                  `json:"expiry_date"`
	Transitional     bool                    `json:"transitional"`
	Type             string                  `json:"type"`
	VersionCreatedAt string                  `json:"version_created_at"`
	Value            certificateVersionValue `json:"value"`
}

type certificateVersionValue struct {
	Ca          string `json:"ca"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

var _ = Describe("Certificates Test", func() {
	Describe("getting a certificate", func() {
		It("should get a certificate with versions", func() {
			certName := "/" + GenerateUniqueCredentialName()
			RunCommand("generate", "-n", certName, "-t", "certificate", "-c", certName, "--is-ca", "--self-sign")

			session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+certName)
			Eventually(session).Should(Exit(0))

			var parsedJson getCertificatesResponse

			err := json.Unmarshal(session.Out.Contents(), &parsedJson)
			Expect(err).ToNot(HaveOccurred())

			Expect(parsedJson.Certificates).To(HaveLen(1))
			Expect(parsedJson.Certificates[0].Versions).To(HaveLen(1))
			Expect(parsedJson.Certificates[0].Versions[0].Id).ToNot(BeEmpty())
			Expect(parsedJson.Certificates[0].Versions[0].ExpiryDate).ToNot(BeEmpty())
			Expect(parsedJson.Certificates[0].Versions[0].Transitional).To(BeFalse())
		})

		Context("when certificate is self-signed", func() {
			It("should get a certificate with signed-by itself", func() {
				certName := "/" + GenerateUniqueCredentialName()
				RunCommand("generate", "-n", certName, "-t", "certificate", "-c", certName, "--is-ca", "--self-sign")

				session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+certName)
				Eventually(session).Should(Exit(0))

				var parsedJson getCertificatesResponse

				err := json.Unmarshal(session.Out.Contents(), &parsedJson)
				Expect(err).ToNot(HaveOccurred())

				Expect(parsedJson.Certificates).To(HaveLen(1))
				Expect(parsedJson.Certificates[0].SignedBy).To(Equal(certName))
			})
		})

		Context("when certificate is signed by a ca", func() {
			It("should get a certificate with signed-by the ca", func() {
				caName := "/" + GenerateUniqueCredentialName()
				certName := "/" + GenerateUniqueCredentialName()
				RunCommand("generate", "-n", caName, "-t", "certificate", "-c", caName, "--is-ca", "--self-sign")
				RunCommand("generate", "-n", certName, "-t", "certificate", "-c", certName, "--ca", caName)

				session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+certName)
				Eventually(session).Should(Exit(0))

				var parsedJson getCertificatesResponse

				err := json.Unmarshal(session.Out.Contents(), &parsedJson)
				Expect(err).ToNot(HaveOccurred())

				Expect(parsedJson.Certificates).To(HaveLen(1))
				Expect(parsedJson.Certificates[0].SignedBy).To(Equal(caName))
			})
		})

		Context("when self-signed certificate is set", func() {
			It("should get a certificate with signed-by itself", func() {
				certName := "/" + GenerateUniqueCredentialName()
				RunCommand("set", "-n", certName, "-t", "certificate", "-c", VALID_CERTIFICATE_CA)

				session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+certName)
				Eventually(session).Should(Exit(0))

				var parsedJson getCertificatesResponse

				err := json.Unmarshal(session.Out.Contents(), &parsedJson)
				Expect(err).ToNot(HaveOccurred())

				Expect(parsedJson.Certificates).To(HaveLen(1))
				Expect(parsedJson.Certificates[0].SignedBy).To(Equal(certName))
			})

			Context("when a certificate signs other certificates", func() {
				It("should return the name of the certificates that it directly signed", func() {
					caName := "/" + GenerateUniqueCredentialName()
					intermediateName1 := "/" + GenerateUniqueCredentialName()
					intermediateName2 := "/" + GenerateUniqueCredentialName()
					leafName1 := "/" + GenerateUniqueCredentialName()
					leafName2 := "/" + GenerateUniqueCredentialName()
					RunCommand("generate", "-n", caName, "-t", "certificate", "-c", caName, "--is-ca", "--self-sign")
					RunCommand("generate", "-n", intermediateName1, "-t", "certificate", "-c", intermediateName1, "--ca", caName)
					RunCommand("generate", "-n", intermediateName2, "-t", "certificate", "-c", intermediateName2, "--ca", caName)
					RunCommand("generate", "-n", leafName1, "-t", "certificate", "-c", leafName1, "--ca", intermediateName1)
					RunCommand("generate", "-n", leafName2, "-t", "certificate", "-c", leafName2, "--ca", intermediateName2)

					session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+caName)
					Eventually(session).Should(Exit(0))

					var parsedJson getCertificatesResponse

					err := json.Unmarshal(session.Out.Contents(), &parsedJson)
					Expect(err).ToNot(HaveOccurred())

					Expect(parsedJson.Certificates).To(HaveLen(1))
					Expect(parsedJson.Certificates[0].Signs).To(ConsistOf(intermediateName1, intermediateName2))
				})
			})
		})

		Context("when non-self-signed certificate is set", func() {
			It("should get a certificate with empty signed-by", func() {
				certName := "/" + GenerateUniqueCredentialName()
				RunCommand("set", "-n", certName, "-t", "certificate", "-c", VALID_CERTIFICATE)

				session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+certName)
				Eventually(session).Should(Exit(0))

				var parsedJson getCertificatesResponse

				err := json.Unmarshal(session.Out.Contents(), &parsedJson)
				Expect(err).ToNot(HaveOccurred())

				Expect(parsedJson.Certificates).To(HaveLen(1))
				Expect(parsedJson.Certificates[0].SignedBy).To(BeEmpty())
			})

			Context("when intermediate ca is set with a ca not in credhub", func() {
				It("should get a certificate with empty signed-by", func() {
					certName := "/" + GenerateUniqueCredentialName()
					RunCommand("set", "-n", certName, "-t", "certificate", "-c", VALID_INTERMEDIATE_CA, "-r", VALID_INTERMEDIATE_CA_ROOT_CA)

					session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates?name="+certName)
					Eventually(session).Should(Exit(0))

					var parsedJson getCertificatesResponse

					err := json.Unmarshal(session.Out.Contents(), &parsedJson)
					Expect(err).ToNot(HaveOccurred())

					Expect(parsedJson.Certificates).To(HaveLen(1))
					Expect(parsedJson.Certificates[0].SignedBy).To(BeEmpty())
				})
			})
		})
	})

	Describe("getting all certificates", func() {
		It("should get all certificates with versions", func() {
			cert1Name := "/" + GenerateUniqueCredentialName()
			cert2Name := "/" + GenerateUniqueCredentialName()
			RunCommand("set", "-n", cert1Name, "-t", "certificate", "--certificate="+VALID_CERTIFICATE, "--private="+VALID_CERTIFICATE_PRIVATE_KEY, "--root="+VALID_CERTIFICATE_CA)
			RunCommand("generate", "-n", cert1Name, "-t", "certificate", "-c", cert1Name, "--is-ca", "--self-sign")
			RunCommand("generate", "-n", cert1Name, "-t", "certificate", "-c", cert1Name, "--is-ca", "--self-sign")
			RunCommand("generate", "-n", cert2Name, "-t", "certificate", "-c", cert2Name, "--is-ca", "--self-sign")
			RunCommand("generate", "-n", cert2Name, "-t", "certificate", "-c", cert2Name, "--is-ca", "--self-sign")

			session := RunCommand("curl", "-X", "GET", "-p", "api/v1/certificates")
			Eventually(session).Should(Exit(0))

			var parsedJson getCertificatesResponse
			err := json.Unmarshal(session.Out.Contents(), &parsedJson)
			Expect(err).ToNot(HaveOccurred())

			cert1Versions := findCert(cert1Name, parsedJson)
			Expect(cert1Versions).NotTo(BeNil())

			cert2Versions := findCert(cert2Name, parsedJson)
			Expect(cert2Versions).NotTo(BeNil())

			Expect(cert1Versions).To(HaveLen(3))
			Expect(cert1Versions[0].Id).ToNot(BeEmpty())
			Expect(cert1Versions[0].ExpiryDate).ToNot(BeEmpty())
			Expect(cert1Versions[0].Transitional).To(BeFalse())
			Expect(cert1Versions[0].Generated).To(BeTrue())
			Expect(cert1Versions[2].Generated).To(BeFalse())
			Expect(cert2Versions).To(HaveLen(2))
			Expect(cert2Versions[0].Id).ToNot(BeEmpty())
			Expect(cert2Versions[0].ExpiryDate).ToNot(BeEmpty())
			Expect(cert2Versions[0].Transitional).To(BeFalse())
		})
	})

	Describe("finding a certificate by ExpiryDate", func() {
		var certificateName string

		BeforeEach(func() {
			certificateName = "/" + GenerateUniqueCredentialName()
			RunCommand("generate", "-n", certificateName, "-t", "certificate", "-d", "15", "-c", certificateName, "--is-ca", "--self-sign")
			RunCommand("generate", "-n", certificateName, "-t", "certificate", "-d", "32", "-c", certificateName, "--is-ca", "--self-sign")
		})

		AfterEach(func() {
			RunCommand("delete", "-n", certificateName)
		})

		Context("using name-like", func() {
			Context("when the latest version of the certificate is not set to expire", func() {
				It("should return no certificates", func() {
					curlCommand := fmt.Sprintf("api/v1/data?name-like=%s&expires-within-days=30", certificateName)
					session := RunCommand("curl", "-X", "GET", "-p", curlCommand)
					Eventually(session).Should(Exit(0))

					stdOut := string(session.Out.Contents())
					Expect(stdOut).To(ContainSubstring(`"credentials": []`))
				})
			})

			Context("when the latest version of the certificate is set to expire", func() {
				It("should return only the latest version of the certificate", func() {
					curlCommand := fmt.Sprintf("api/v1/data?name-like=%s&expires-within-days=50", certificateName)
					session := RunCommand("curl", "-X", "GET", "-p", curlCommand)
					Eventually(session).Should(Exit(0))

					stdOut := string(session.Out.Contents())
					Expect(stdOut).To(ContainSubstring(`"name": "` + certificateName + `"`))
				})
			})
		})

		Context("using path", func() {
			Context("when the latest version of the certificate is not set to expire", func() {
				It("should return no certificates", func() {
					curlCommand := "api/v1/data?path=/&expires-within-days=30"
					session := RunCommand("curl", "-X", "GET", "-p", curlCommand)
					Eventually(session).Should(Exit(0))

					stdOut := string(session.Out.Contents())
					Expect(stdOut).To(ContainSubstring(`"credentials":`))
					Expect(stdOut).NotTo(ContainSubstring(`"name": "` + certificateName + `"`))
				})
			})

			Context("when the latest version of the certificate is set to expire", func() {
				It("should return only the latest version of the certificate", func() {
					curlCommand := "api/v1/data?path=/&expires-within-days=50"
					session := RunCommand("curl", "-X", "GET", "-p", curlCommand)
					Eventually(session).Should(Exit(0))

					stdOut := string(session.Out.Contents())
					Expect(stdOut).To(ContainSubstring(`"name": "` + certificateName + `"`))
				})
			})
		})
	})

	Describe("setting a certificate", func() {
		Context("when private key format is PKCS1", func() {
			It("should be able to set a certificate", func() {
				name := GenerateUniqueCredentialName()
				RunCommand("set", "-n", name, "-t", "certificate", "--certificate="+VALID_CERTIFICATE, "--private="+VALID_CERTIFICATE_PRIVATE_KEY, "--root="+VALID_CERTIFICATE_CA)
				session := RunCommand("get", "-n", name)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				expectedCertValue := CertificateValue{
					Ca:          VALID_CERTIFICATE_CA,
					Certificate: VALID_CERTIFICATE,
					PrivateKey:  VALID_CERTIFICATE_PRIVATE_KEY,
				}

				Eventually(session).Should(Exit(0))

				actualCert := Certificate{}
				err := yaml.Unmarshal([]byte(stdOut), &actualCert)
				Expect(err).To(BeNil())
				Expect(actualCert.Name).To(Equal("/" + name))
				Expect(actualCert.Value).To(Equal(expectedCertValue))

			})

			It("should require a certificate type", func() {
				session := RunCommand("set", "-n", GenerateUniqueCredentialName(), "-t", "certificate")
				Eventually(session).Should(Exit(1))
				Expect(session.Err.Contents()).To(MatchRegexp(".*At least one certificate attribute must be set. Please validate your input and retry your request."))
			})

			It("should allow you to set a certificate with a named CA", func() {
				caName := GenerateUniqueCredentialName()
				certName := GenerateUniqueCredentialName()
				RunCommand("set", "-n", caName, "-t", "certificate", "-c", VALID_CERTIFICATE_CA)
				session := RunCommand("get", "-n", caName)
				Eventually(session).Should(Exit(0))
				stdOut := string(session.Out.Contents())

				caCert := Certificate{}
				err := yaml.Unmarshal([]byte(stdOut), &caCert)
				Expect(err).To(BeNil())

				RunCommand("set", "-n", certName, "-t", "certificate", "--certificate="+VALID_CERTIFICATE, "--private="+VALID_CERTIFICATE_PRIVATE_KEY, "--ca-name", caName)
				session = RunCommand("get", "-n", certName)
				Eventually(session).Should(Exit(0))
				stdOut = string(session.Out.Contents())
				cert := Certificate{}
				err = yaml.Unmarshal([]byte(stdOut), &cert)
				Expect(err).To(BeNil())
				Expect(cert.Value.Ca).To(Equal(caCert.Value.Certificate))
				Expect(cert.Name).To(Equal("/" + certName))
				Expect(cert.Value.PrivateKey).To(Equal(VALID_CERTIFICATE_PRIVATE_KEY))
				Expect(cert.Value.Certificate).To(Equal(VALID_CERTIFICATE))

			})
		})

		Context("when private key format is PKCS8", func() {
			Context("and is RSA formatted", func() {
				It("should store certificate in database", func() {
					name := GenerateUniqueCredentialName()
					RunCommand("set", "-n", name, "-t", "certificate", "--certificate="+OTHER_VALID_CERTIFICATE, "--private="+OTHER_VALID_PRIVATE_KEY_PKCS8)
					session := RunCommand("get", "-n", name)
					stdOut := string(session.Out.Contents())

					expectedCertValue := CertificateValue{
						Certificate: OTHER_VALID_CERTIFICATE,
						PrivateKey:  OTHER_VALID_PRIVATE_KEY_PKCS8,
					}

					Eventually(session).Should(Exit(0))

					actualCert := Certificate{}
					err := yaml.Unmarshal([]byte(stdOut), &actualCert)
					Expect(err).To(BeNil())
					Expect(actualCert.Name).To(Equal("/" + name))
					Expect(actualCert.Value).To(Equal(expectedCertValue))
				})
			})

			Context("and is not RSA formatted", func() {
				It("should return an error", func() {
					name := GenerateUniqueCredentialName()
					session := RunCommand("set", "-n", name, "-t", "certificate", "--certificate="+OTHER_VALID_CERTIFICATE, "--private="+EC_PRIVATE_KEY)
					stdErr := strings.TrimSpace(string(session.Err.Contents()))

					Eventually(session).Should(Exit(1))
					Expect(stdErr).To(Equal("Private key is malformed. Key file does not contain an RSA private key"))
				})
			})
			Context("and is encrypted", func() {
				It("should return an error", func() {
					name := GenerateUniqueCredentialName()
					session := RunCommand("set", "-n", name, "-t", "certificate", "--certificate="+OTHER_VALID_CERTIFICATE, "--private="+OTHER_PRIVATE_KEY_PKCS8_ENCRYPTED)
					stdErr := strings.TrimSpace(string(session.Err.Contents()))

					Eventually(session).Should(Exit(1))
					Expect(stdErr).To(Equal("Private key is malformed. Key file is not in PKCS#1 or unencrypted PKCS#8 format"))
				})
			})
		})

	})

	Describe("CAs and Certificates", func() {
		Describe("certificate chains", func() {
			It("should build the chain with an intermediate CA", func() {
				rootCaName := GenerateUniqueCredentialName()
				intermediateCaName := GenerateUniqueCredentialName()
				leafCertificateName := GenerateUniqueCredentialName()

				RunCommand("generate", "-n", rootCaName, "-t", "certificate", "-c", rootCaName, "--is-ca", "--self-sign")
				session := RunCommand("get", "-n", rootCaName)
				rootCert := CertFromPem(string(session.Out.Contents()), false)
				Expect(rootCert.Subject.CommonName).To(Equal(rootCaName))
				Expect(rootCert.Issuer.CommonName).To(Equal(rootCaName))
				Expect(rootCert.IsCA).To(Equal(true))
				Expect(len(rootCert.SubjectKeyId)).ToNot(Equal(0))

				RunCommand("generate", "-n", intermediateCaName, "-t", "certificate", "-c", intermediateCaName, "--is-ca", "--ca", rootCaName)
				session = RunCommand("get", "-n", intermediateCaName)
				intermediateCert := CertFromPem(string(session.Out.Contents()), false)
				Expect(intermediateCert.Subject.CommonName).To(Equal(intermediateCaName))
				Expect(intermediateCert.Issuer.CommonName).To(Equal(rootCert.Subject.CommonName))
				Expect(intermediateCert.IsCA).To(Equal(true))

				RunCommand("generate", "-n", leafCertificateName, "-t", "certificate", "-c", leafCertificateName, "--ca", intermediateCaName)
				session = RunCommand("get", "-n", leafCertificateName)
				leafCert := CertFromPem(string(session.Out.Contents()), false)

				caCerts, err := ioutil.TempFile("", "credhubTestCerts")
				Expect(err).ToNot(HaveOccurred())
				pem.Encode(caCerts, &pem.Block{Type: "CERTIFICATE", Bytes: rootCert.Raw})
				pem.Encode(caCerts, &pem.Block{Type: "CERTIFICATE", Bytes: intermediateCert.Raw})
				caCerts.Close()

				leafFile, err := ioutil.TempFile("", "leafCert")
				Expect(err).ToNot(HaveOccurred())
				pem.Encode(leafFile, &pem.Block{Type: "CERTIFICATE", Bytes: leafCert.Raw})
				leafFile.Close()

				cmd := exec.Command("openssl", "verify", "-CAfile", caCerts.Name(), leafFile.Name())
				session, err = Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(gbytes.Say(fmt.Sprintf("%s: OK", leafFile.Name())))
				Expect(leafCert.Subject.CommonName).To(Equal(leafCertificateName))
				Expect(leafCert.Issuer.CommonName).To(Equal(intermediateCert.Subject.CommonName))
				Expect(leafCert.IsCA).To(Equal(false))
			})
		})

		It("should return multiple CAs if the concatenate CA flag is set", func() {
			caName := "/" + GenerateUniqueCredentialName()
			certName := "/" + GenerateUniqueCredentialName()

			session := RunCommand("generate", "-n", caName, "-t", "certificate", "-c", caName, "--is-ca", "--self-sign")
			Expect(session).To(Exit(0))

			session = RunCommand("generate", "-n", certName, "-t", "certificate", "-c", certName, "--ca", caName)
			Expect(session).To(Exit(0))

			session = RunCommand("curl", "-p", fmt.Sprintf("/api/v1/certificates?name=%s", caName))
			Expect(session).To(Exit(0))

			var response getCertificatesResponse
			err := json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())

			session = RunCommand("curl",
				"-p", fmt.Sprintf("/api/v1/certificates/%s/regenerate", response.Certificates[0].Id),
				"-X", "POST",
				"-d", `{"set_as_transitional": "true"}`,
			)
			Expect(session).To(Exit(0))

			session = RunCommand("get", "-n", certName, "-k", "ca")
			Expect(session).To(Exit(0))
			stdOut := string(session.Out.Contents())

			re := regexp.MustCompile("BEGIN CERTIFICATE")
			certificates := re.FindAllString(stdOut, -1)
			if cfg.ConcatenateCas {
				Expect(certificates).To(HaveLen(2))
			} else {
				Expect(certificates).To(HaveLen(1))
			}

			session = RunCommand("curl", "-p", fmt.Sprintf("/api/v1/certificates?name=%s", certName))
			Expect(session).To(Exit(0))
			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())

			var versionResponse getCertificateVersionsResponse
			session = RunCommand("curl", "-p", fmt.Sprintf("/api/v1/certificates/%s/versions", response.Certificates[0].Id))
			Expect(session).To(Exit(0))
			err = json.Unmarshal(session.Out.Contents(), &versionResponse)
			Expect(err).NotTo(HaveOccurred())

			ca := versionResponse[0].Value.Ca
			certificates = re.FindAllString(ca, -1)
			if cfg.ConcatenateCas {
				Expect(certificates).To(HaveLen(2))
			} else {
				Expect(certificates).To(HaveLen(1))
			}
		})

		It("should generate a ca when using the --is-ca flag", func() {
			certificateId := GenerateUniqueCredentialName()
			certificateAuthorityId := GenerateUniqueCredentialName()

			By("generating the CA", func() {
				RunCommand("generate", "-n", certificateAuthorityId, "-t", "certificate", "--common-name", certificateAuthorityId, "--is-ca")
				session := RunCommand("get", "-n", certificateAuthorityId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(MatchRegexp(`certificate: |\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
				cert := CertFromPem(stdOut, false)
				Expect(cert.Subject.CommonName).To(Equal(certificateAuthorityId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId)) // self-signed
				Expect(cert.IsCA).To(Equal(true))
			})

			By("getting the CA", func() {
				session := RunCommand("get", "-n", certificateAuthorityId)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				cert := CertFromPem(stdOut, false)
				Expect(cert.Subject.CommonName).To(Equal(certificateAuthorityId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId)) // self-signed
				Expect(cert.IsCA).To(Equal(true))
			})

			By("generating and signing the certificate", func() {
				RunCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_signing", "-g", "digital_signature", "-a", "example.com", "-k", "3072", "-d", "90")
				session := RunCommand("get", "-n", certificateId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(MatchRegexp(`certificate: |\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))
				cert := CertFromPem(stdOut, false)
				ca := CertFromPem(stdOut, true)

				Expect(cert.AuthorityKeyId).To(Equal(ca.SubjectKeyId))

				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId))
				Expect(ca.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)).To(BeNil()) // signed by ca
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.IsCA).To(Equal(false))
				Expect(cert.NotAfter.Sub(cert.NotBefore).Hours()).To(Equal(90 * 24.0))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(3072))
				Expect(cert.DNSNames).To(Equal([]string{"example.com"}))
			})

			By("getting the certificate", func() {
				session := RunCommand("get", "-n", certificateId)
				Eventually(session).Should(Exit(0))
			})

			By("regenerating the certificate", func() {
				RunCommand("regenerate", "-n", certificateId)
				session := RunCommand("get", "-n", certificateId)
				Eventually(session).Should(Exit(0))
				stdOut := string(session.Out.Contents())
				cert := CertFromPem(stdOut, false)
				ca := CertFromPem(stdOut, true)
				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId))
				Expect(ca.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)).To(BeNil()) // signed by ca
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.IsCA).To(Equal(false))
				Expect(cert.NotAfter.Sub(cert.NotBefore).Hours()).To(Equal(90 * 24.0))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(3072))
				Expect(cert.DNSNames).To(Equal([]string{"example.com"}))
			})
		})

		It("should be able to generate a self-signed certificate", func() {
			certificateId := GenerateUniqueCredentialName()
			initialCertificate := ""
			initialPrivateKey := ""

			By("generating the certificate", func() {
				RunCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--self-sign", "-e", "email_protection", "-g", "digital_signature", "-a", "example.com", "-k", "3072", "-d", "90")
				session := RunCommand("get", "-n", certificateId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`type: certificate`))
				Expect(stdOut).To(MatchRegexp(`certificate: |\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`private_key: |\s+-----BEGIN RSA PRIVATE KEY-----`))

				initialCertificate = stdOut[strings.Index(stdOut, "-----BEGIN CERTIFICATE-----"):strings.Index(stdOut, "-----END CERTIFICATE-----")]
				initialPrivateKey = stdOut[strings.Index(stdOut, "-----BEGIN RSA PRIVATE KEY-----"):strings.Index(stdOut, "-----END RSA PRIVATE KEY-----")]

				cert := CertFromPem(stdOut, false)
				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateId))                                                  // self-signed
				Expect(cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)).To(BeNil()) // signed by self
				Expect(cert.IsCA).To(Equal(false))
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection}))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.NotAfter.Sub(cert.NotBefore).Hours()).To(Equal(90 * 24.0))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(3072))
				Expect(cert.DNSNames).To(Equal([]string{"example.com"}))
			})

			By("getting the certificate", func() {
				session := RunCommand("get", "-n", certificateId)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(MatchRegexp(`certificate: |\s+-----BEGIN CERTIFICATE-----`))
			})

			By("regenerating the certificate", func() {
				RunCommand("regenerate", "-n", certificateId)
				session := RunCommand("get", "-n", certificateId)
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				cert := CertFromPem(stdOut, false)
				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateId))                                                  // self-signed
				Expect(cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)).To(BeNil()) // signed by self
				Expect(cert.IsCA).To(Equal(false))
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection}))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.NotAfter.Sub(cert.NotBefore).Hours()).To(Equal(90 * 24.0))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(3072))
				Expect(cert.DNSNames).To(Equal([]string{"example.com"}))

				Expect(stdOut).NotTo(ContainSubstring(initialCertificate))
				Expect(stdOut).NotTo(ContainSubstring(initialPrivateKey))
			})
		})

		It("should error gracefully when supplying an invalid extended key usage name", func() {
			certificateAuthorityId := GenerateUniqueCredentialName()
			certificateId := certificateAuthorityId + "1"
			RunCommand("generate", "-n", certificateAuthorityId, "-t certificate", "--common-name", certificateAuthorityId, "--is-ca")
			session := RunCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_sinning")
			stdErr := string(session.Err.Contents())

			Eventually(session).Should(Exit(1))
			Expect(stdErr).To(MatchRegexp(`The provided extended key usage 'code_sinning' is not supported. Valid values include 'client_auth', 'server_auth', 'code_signing', 'email_protection' and 'timestamping'.`))
		})

		It("should error gracefully when supplying an invalid key usage name", func() {
			certificateAuthorityId := GenerateUniqueCredentialName()
			certificateId := certificateAuthorityId + "1"
			RunCommand("generate", "-n", certificateAuthorityId, "-t certificate", "--common-name", certificateAuthorityId, "--is-ca")
			session := RunCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-g", "digital_sinnature")
			stdErr := string(session.Err.Contents())

			Eventually(session).Should(Exit(1))
			Expect(stdErr).To(MatchRegexp(`The provided key usage 'digital_sinnature' is not supported. Valid values include 'digital_signature', 'non_repudiation', 'key_encipherment', 'data_encipherment', 'key_agreement', 'key_cert_sign', 'crl_sign', 'encipher_only' and 'decipher_only'.`))
		})
	})
})

func findCert(certName string, certs getCertificatesResponse) []certificateVersion {
	for _, certificate := range certs.Certificates {
		if certName == certificate.Name {
			return certificate.Versions
		}
	}
	return nil
}

// https://golang.org/pkg/crypto/x509/#Certificate
// prefix should be "Certificate" or "Ca"
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
