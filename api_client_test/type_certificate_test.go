package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

var _ = Describe("Certificate Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-certificate")

		generateCert := generate.Certificate{
			CommonName: "example.com",
			SelfSign:   true,
		}

		setCert := values.Certificate{
			Ca:          test_helpers.VALID_CERTIFICATE_CA,
			Certificate: test_helpers.VALID_CERTIFICATE,
			PrivateKey:  test_helpers.VALID_CERTIFICATE_PRIVATE_KEY,
		}

		By("generate a certificate with path " + name)
		certificate, err := credhubClient.GenerateCertificate(name, generateCert, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value.Certificate).ToNot(BeEmpty())
		Expect(certificate.Value.PrivateKey).ToNot(BeEmpty())
		firstGeneratedCertificate := certificate.Value

		By("generate the certificate again without overwrite returns same certificate")
		certificate, err = credhubClient.GenerateCertificate(name, generateCert, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(firstGeneratedCertificate))

		By("setting the certificate again without overwrite returns same certificate")
		certificate, err = credhubClient.SetCertificate(name, setCert, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(firstGeneratedCertificate))

		By("overwriting the certificate with generate")
		certificate, err = credhubClient.GenerateCertificate(name, generateCert, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).ToNot(Equal(firstGeneratedCertificate))

		By("overwriting the certificate with a provded CA name")
		generateCA := generate.Certificate{
			CommonName: "some-ca",
			SelfSign:   true,
			IsCA:       true,
		}

		ca, err := credhubClient.GenerateCertificate("/test-ca", generateCA, credhub.Overwrite)

		setCert = values.Certificate{
			CaName:      "/test-ca",
			Certificate: test_helpers.VALID_CERTIFICATE,
			PrivateKey:  test_helpers.VALID_CERTIFICATE_PRIVATE_KEY,
		}

		certificate, err = credhubClient.SetCertificate(name, setCert, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value.Ca).To(Equal(ca.Value.Certificate))

		By("getting the certificate")
		certificate, err = credhubClient.GetLatestCertificate(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value.Ca).To(Equal(ca.Value.Certificate))

		By("deleting the certificate")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		err = credhubClient.Delete("/test-ca")
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestCertificate(name)
		Expect(err).To(HaveOccurred())
	})
})
