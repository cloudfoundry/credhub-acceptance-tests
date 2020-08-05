package certs_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"time"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers/certs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {

	const THRESHOLD = 2 * time.Minute

	Describe("GenerateSelfSigned", func() {
		It("generates a valid self-signed certificate", func() {
			certBytes, keyBytes, err := GenerateSelfSigned(CertOptions{})
			Expect(err).NotTo(HaveOccurred())

			cert := parseCert(certBytes, keyBytes)
			Expect(cert).To(BeValidSelfSignedCert())
			Expect(cert.Subject.CommonName).To(BeEmpty())
			Expect(cert.Subject.OrganizationalUnit).To(BeNil())
			Expect(cert.ExtKeyUsage).To(ConsistOf(x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth))
			Expect(cert.SerialNumber).NotTo(BeNil())
		})

		Context("when a subject is provided", func() {
			It("generates a self-signed certificate with the given subject", func() {
				commonName := "some-common-name"
				organizationalUnit := "some-organizational-unit"

				options := CertOptions{CommonName: commonName, OrganizationalUnit: organizationalUnit}
				certBytes, keyBytes, err := GenerateSelfSigned(options)
				Expect(err).NotTo(HaveOccurred())

				cert := parseCert(certBytes, keyBytes)
				Expect(cert).To(BeValidSelfSignedCert())
				Expect(cert.Subject.CommonName).To(Equal(commonName))
				Expect(cert.Subject.OrganizationalUnit).To(Equal([]string{organizationalUnit}))
			})
		})

		Context("when IsCA is true", func() {
			It("generates a self-signed CA certificate", func() {
				options := CertOptions{IsCA: true}
				certBytes, keyBytes, err := GenerateSelfSigned(options)
				Expect(err).NotTo(HaveOccurred())

				cert := parseCert(certBytes, keyBytes)
				Expect(cert).To(BeValidSelfSignedCert())
				Expect(cert.IsCA).To(BeTrue())
				Expect(int(cert.KeyUsage & x509.KeyUsageCertSign)).NotTo(Equal(0))
			})
		})

		Describe("validity dates", func() {
			It("generates a certificate valid immediately and expiring in 30 days", func() {
				now := time.Now().UTC()
				thirtyDaysFromNow := now.Add(time.Hour * 24 * 30)

				options := CertOptions{}
				certBytes, keyBytes, err := GenerateSelfSigned(options)
				Expect(err).NotTo(HaveOccurred())

				cert := parseCert(certBytes, keyBytes)
				Expect(cert).To(BeValidSelfSignedCert())
				Expect(cert.NotBefore).To(BeTemporally("~", now, THRESHOLD))
				Expect(cert.NotAfter).To(BeTemporally("~", thirtyDaysFromNow, THRESHOLD))
			})

			Context("given only NotBefore", func() {
				It("generates a certificate expiring 30 days after NotBefore", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10 * -1)
					thirtyDaysLater := notBefore.Add(time.Hour * 24 * 30)

					options := CertOptions{NotBefore: notBefore}
					certBytes, keyBytes, err := GenerateSelfSigned(options)
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(BeValidSelfSignedCert())
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", thirtyDaysLater, THRESHOLD))
				})
			})

			Context("given only NotAfter", func() {
				It("generates a certificate valid immediately and expiring on NotAfter", func() {
					now := time.Now().UTC()
					notAfter := now.Add(time.Hour * 24 * 10)

					options := CertOptions{NotAfter: notAfter}
					certBytes, keyBytes, err := GenerateSelfSigned(options)
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(BeValidSelfSignedCert())
					Expect(cert.NotBefore).To(BeTemporally("~", now, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", notAfter, THRESHOLD))
				})
			})

			Context("given NotBefore and NotAfter", func() {
				It("generates a certificate using both dates", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10 * -1)
					notAfter := time.Now().UTC().Add(time.Hour * 24 * 20)

					options := CertOptions{NotBefore: notBefore, NotAfter: notAfter}
					certBytes, keyBytes, err := GenerateSelfSigned(options)
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(BeValidSelfSignedCert())
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", notAfter, THRESHOLD))
				})
			})

			Context("given a NotBefore in the future", func() {
				It("generates a certificate that is not yet valid", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10)

					options := CertOptions{NotBefore: notBefore}
					certBytes, keyBytes, err := GenerateSelfSigned(options)
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(FailCertValidationWithMessage("x509: certificate has expired or is not yet valid"))
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
				})
			})

			Context("given a NotAfter in the past", func() {
				It("generates a certificate that is expired", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10 * -1)
					notAfter := time.Now().UTC().Add(time.Hour * 24 * 5 * -1)

					options := CertOptions{NotBefore: notBefore, NotAfter: notAfter}
					certBytes, keyBytes, err := GenerateSelfSigned(options)
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(FailCertValidationWithMessage("x509: certificate has expired or is not yet valid"))
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", notAfter, THRESHOLD))
				})
			})

			Context("when NotAfter is earlier than NotBefore", func() {
				It("returns an error", func() {
					notAfter := time.Now().UTC()
					notBefore := notAfter.Add(time.Hour * 24)

					options := CertOptions{NotBefore: notBefore, NotAfter: notAfter}
					_, _, err := GenerateSelfSigned(options)
					Expect(err).To(MatchError(MatchRegexp(`NotBefore (.*) must be earlier than NotAfter (.*)`)))
				})
			})
		})
	})

	Describe("GenerateSigned", func() {
		It("generates a valid certificate signed by the given CA", func() {
			certBytes, keyBytes, err := GenerateSigned(CertOptions{}, []byte(CaCert), []byte(CaKey))
			Expect(err).NotTo(HaveOccurred())

			cert := parseCert(certBytes, keyBytes)
			Expect(cert).To(BeValidCertSignedBy([]byte(CaCert)))
			Expect(cert.Subject.CommonName).To(BeEmpty())
			Expect(cert.Subject.OrganizationalUnit).To(BeNil())
			Expect(cert.ExtKeyUsage).To(ConsistOf(x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth))
			Expect(cert.SerialNumber).NotTo(BeNil())
		})

		Context("when no subject is provided", func() {
			It("generates a certificate with the given subject signed by the given CA", func() {
				commonName := "some-common-name"
				organizationalUnit := "some-organizational-unit"

				options := CertOptions{CommonName: commonName, OrganizationalUnit: organizationalUnit}
				certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
				Expect(err).NotTo(HaveOccurred())

				cert := parseCert(certBytes, keyBytes)
				Expect(cert).To(BeValidSelfSignedCert())
				Expect(cert.Subject.CommonName).To(Equal(commonName))
				Expect(cert.Subject.OrganizationalUnit).To(Equal([]string{organizationalUnit}))
			})
		})

		Context("when IsCA is true", func() {
			It("generates a CA certificate signed by the given CA", func() {
				options := CertOptions{IsCA: true}
				certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
				Expect(err).NotTo(HaveOccurred())

				cert := parseCert(certBytes, keyBytes)
				Expect(cert).To(BeValidCertSignedBy([]byte(CaCert)))
				Expect(cert.IsCA).To(BeTrue())
				Expect(int(cert.KeyUsage & x509.KeyUsageCertSign)).NotTo(Equal(0))
			})
		})

		Describe("validity dates", func() {
			It("generates a certificate valid immediately and expiring in 30 days", func() {
				now := time.Now().UTC()
				thirtyDaysFromNow := now.Add(time.Hour * 24 * 30)

				options := CertOptions{}
				certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
				Expect(err).NotTo(HaveOccurred())

				cert := parseCert(certBytes, keyBytes)
				Expect(cert.NotBefore).To(BeTemporally("~", now, THRESHOLD))
				Expect(cert.NotAfter).To(BeTemporally("~", thirtyDaysFromNow, THRESHOLD))
			})

			Context("given only NotBefore", func() {
				It("generates a certificate expiring 30 days after NotBefore", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10 * -1)
					thirtyDaysLater := notBefore.Add(time.Hour * 24 * 30)

					options := CertOptions{NotBefore: notBefore}
					certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", thirtyDaysLater, THRESHOLD))
				})
			})

			Context("given only NotAfter", func() {
				It("generates a certificate valid immediately and expiring on NotAfter", func() {
					now := time.Now().UTC()
					notAfter := now.Add(time.Hour * 24 * 10)

					options := CertOptions{NotAfter: notAfter}
					certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert.NotBefore).To(BeTemporally("~", now, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", notAfter, THRESHOLD))
				})
			})

			Context("given NotBefore and NotAfter", func() {
				It("generates a certificate using both dates", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10 * -1)
					notAfter := time.Now().UTC().Add(time.Hour * 24 * 20)

					options := CertOptions{NotBefore: notBefore, NotAfter: notAfter}
					certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", notAfter, THRESHOLD))
				})
			})

			Context("given a NotBefore in the future", func() {
				It("generates a certificate that is not yet valid", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10)

					options := CertOptions{NotBefore: notBefore}
					certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(FailCertValidationWithMessage("x509: certificate has expired or is not yet valid"))
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
				})
			})

			Context("given a NotAfter in the past", func() {
				It("generates a certificate that is expired", func() {
					notBefore := time.Now().UTC().Add(time.Hour * 24 * 10 * -1)
					notAfter := time.Now().UTC().Add(time.Hour * 24 * 5 * -1)

					options := CertOptions{NotBefore: notBefore, NotAfter: notAfter}
					certBytes, keyBytes, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
					Expect(err).NotTo(HaveOccurred())

					cert := parseCert(certBytes, keyBytes)
					Expect(cert).To(FailCertValidationWithMessage("x509: certificate has expired or is not yet valid"))
					Expect(cert.NotBefore).To(BeTemporally("~", notBefore, THRESHOLD))
					Expect(cert.NotAfter).To(BeTemporally("~", notAfter, THRESHOLD))
				})
			})

			Context("when NotAfter is earlier than NotBefore", func() {
				It("returns an error", func() {
					notAfter := time.Now().UTC()
					notBefore := notAfter.Add(time.Hour * 24)

					options := CertOptions{NotBefore: notBefore, NotAfter: notAfter}
					_, _, err := GenerateSigned(options, []byte(CaCert), []byte(CaKey))
					Expect(err).To(MatchError(MatchRegexp(`NotBefore (.*) must be earlier than NotAfter (.*)`)))
				})
			})
		})

		Context("when the CA cert and key are invalid", func() {
			It("returns an error", func() {
				_, _, err := GenerateSigned(CertOptions{}, nil, nil)
				Expect(err).To(MatchError(ContainSubstring("failed to load CA key pair")))
			})
		})
	})
})

func parseCert(certBytes, keyBytes []byte) *x509.Certificate {
	_, err := tls.X509KeyPair(certBytes, keyBytes)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	block, rest := pem.Decode(certBytes)
	ExpectWithOffset(1, block).NotTo(BeNil())
	ExpectWithOffset(1, rest).To(BeEmpty())

	cert, err := x509.ParseCertificate(block.Bytes)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return cert
}

func verifySignedBy(cert *x509.Certificate, caCertBytes []byte) {
	roots := x509.NewCertPool()
	ExpectWithOffset(1, roots.AppendCertsFromPEM(caCertBytes)).To(BeTrue())
	_, err := cert.Verify(x509.VerifyOptions{Roots: roots})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}
