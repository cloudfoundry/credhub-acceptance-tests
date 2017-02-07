package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"crypto/x509"
	"encoding/pem"
	"regexp"
	"crypto/rsa"
)

var _ = Describe("Certificates Test", func() {
	Describe("setting a certificate", func() {
		It("should be able to set a certificate", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "certificate", "--certificate-string", "iamacertificate")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
			Expect(stdOut).To(MatchRegexp(`Certificate:\s+iamacertificate`))
		})

		It("should require a certificate type", func() {
			session := runCommand("set", "-n", generateUniqueCredentialName(), "-t", "certificate", "--no-overwrite")
			Eventually(session).Should(Exit(1))
			Expect(session.Err.Contents()).To(MatchRegexp(".*At least one certificate type must be set. Please validate your input and retry your request."))
		})
	})

	Describe("CAs and Certificates", func() {
		Describe("certificate chains", func() {
			It("should build the chain with an intermediate CA", func() {
				rootCaName := generateUniqueCredentialName()
				intermediateCaName := generateUniqueCredentialName()
				leafCertificateName := generateUniqueCredentialName()

				session := runCommand("generate", "-n", rootCaName, "-t", "certificate", "-c", rootCaName, "--is-ca", "--self-sign")
				cert := CertFromPem(string(session.Out.Contents()), "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(rootCaName))
				Expect(cert.Issuer.CommonName).To(Equal(rootCaName))
				Expect(cert.IsCA).To(Equal(true))

				session = runCommand("generate", "-n", intermediateCaName, "-t", "certificate", "-c", intermediateCaName, "--is-ca", "--ca", rootCaName)
				cert = CertFromPem(string(session.Out.Contents()), "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(intermediateCaName))
				Expect(cert.Issuer.CommonName).To(Equal(rootCaName))
				Expect(cert.IsCA).To(Equal(true))

				session = runCommand("generate", "-n", leafCertificateName, "-t", "certificate", "-c", leafCertificateName, "--ca", intermediateCaName)
				cert = CertFromPem(string(session.Out.Contents()), "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(leafCertificateName))
				Expect(cert.Issuer.CommonName).To(Equal(intermediateCaName))
				Expect(cert.IsCA).To(Equal(false))
			})
		})

		It("it should generate CA and certificate when using the 'ca-' commands, ", func() {
			certificateAuthorityId := generateUniqueCredentialName()
			certificateId := certificateAuthorityId + "1"

			By("retrieving a CA that doesn't exist yet", func() {
				session := runCommand("ca-get", "-n", certificateAuthorityId)
				stdErr := string(session.Err.Contents())

				Expect(stdErr).To(MatchRegexp(`The request could not be completed because the CA has not been defined. Please set the CA and retry your request.`))
				Eventually(session).Should(Exit(1))
			})

			By("generating the CA", func() {
				session := runCommand("ca-generate", "-n", certificateAuthorityId, "--common-name", certificateAuthorityId)
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+root`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
			})

			By("getting the new CA", func() {
				session := runCommand("ca-get", "-n", certificateAuthorityId)
				Eventually(session).Should(Exit(0))
			})

			By("generating and signing the certificate", func() {
				session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_signing", "-g", "digital_signature")
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
				cert := CertFromPem(stdOut, "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}))
				Expect(cert.IsCA).To(Equal(false))
			})

			By("getting the certificate", func() {
				session := runCommand("get", "-n", certificateId)
				Eventually(session).Should(Exit(0))
			})
		})

		It("should generate a ca when using the --is-ca flag", func() {
			certificateId := generateUniqueCredentialName()
			certificateAuthorityId := generateUniqueCredentialName()

			By("generating the CA", func() {
				session := runCommand("generate", "-n", certificateAuthorityId, "-t", "certificate", "--common-name", certificateAuthorityId, "--is-ca")
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
				cert := CertFromPem(stdOut, "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(certificateAuthorityId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId)) // self-signed
				Expect(cert.IsCA).To(Equal(true))
			})

			By("getting the CA", func() {
				session := runCommand("get", "-n", certificateAuthorityId)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				cert := CertFromPem(stdOut, "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(certificateAuthorityId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateAuthorityId)) // self-signed
				Expect(cert.IsCA).To(Equal(true))
			})

			By("generating and signing the certificate", func() {
				session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_signing", "-g", "digital_signature", "-a", "example.com", "-k", "3072", "-d", "90")
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
				cert := CertFromPem(stdOut, "Certificate")
				ca := CertFromPem(stdOut, "Ca")
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
				session := runCommand("get", "-n", certificateId)
				Eventually(session).Should(Exit(0))
			})

			By("regenerating the certificate", func() {
				session := runCommand("regenerate", "-n", certificateId)
				Eventually(session).Should(Exit(0))
				stdOut := string(session.Out.Contents())
				cert := CertFromPem(stdOut, "Certificate")
				ca := CertFromPem(stdOut, "Ca")
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
			certificateId := generateUniqueCredentialName()
			By("generating the certificate", func() {
				session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--self-sign", "-e", "email_protection", "-g", "digital_signature", "-a", "example.com", "-k", "3072", "-d", "90")
				stdOut := string(session.Out.Contents())

				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).ToNot(MatchRegexp(`Ca:`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))

				cert := CertFromPem(stdOut, "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateId)) // self-signed
				Expect(cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)).To(BeNil()) // signed by self
				Expect(cert.IsCA).To(Equal(false))
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection}))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.NotAfter.Sub(cert.NotBefore).Hours()).To(Equal(90 * 24.0))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(3072))
				Expect(cert.DNSNames).To(Equal([]string{"example.com"}))
			})

			By("getting the certificate", func() {
				session := runCommand("get", "-n", certificateId)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				Expect(stdOut).ToNot(MatchRegexp(`Ca:`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
			})

			By("regenerating the certificate", func() {
				session := runCommand("regenerate", "-n", certificateId)
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				cert := CertFromPem(stdOut, "Certificate")
				Expect(cert.Subject.CommonName).To(Equal(certificateId))
				Expect(cert.Issuer.CommonName).To(Equal(certificateId)) // self-signed
				Expect(cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)).To(BeNil()) // signed by self
				Expect(cert.IsCA).To(Equal(false))
				Expect(cert.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection}))
				Expect(cert.KeyUsage).To(Equal(x509.KeyUsageDigitalSignature))
				Expect(cert.NotAfter.Sub(cert.NotBefore).Hours()).To(Equal(90 * 24.0))
				Expect(cert.PublicKey.(*rsa.PublicKey).N.BitLen()).To(Equal(3072))
				Expect(cert.DNSNames).To(Equal([]string{"example.com"}))
			})
		})

		It("should error gracefully when supplying an invalid extended key usage name", func() {
			certificateAuthorityId := generateUniqueCredentialName()
			certificateId := certificateAuthorityId + "1"
			runCommand("ca-generate", "-n", certificateAuthorityId, "--common-name", certificateAuthorityId)
			session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-e", "code_sinning")
			stdErr := string(session.Err.Contents())

			Eventually(session).Should(Exit(1))
			Expect(stdErr).To(MatchRegexp(`The provided extended key usage 'code_sinning' was not known. Please update this value and retry your request.`))
		})

		It("should error gracefully when supplying an invalid key usage name", func() {
			certificateAuthorityId := generateUniqueCredentialName()
			certificateId := certificateAuthorityId + "1"
			runCommand("ca-generate", "-n", certificateAuthorityId, "--common-name", certificateAuthorityId)
			session := runCommand("generate", "-n", certificateId, "-t", "certificate", "--common-name", certificateId, "--ca", certificateAuthorityId, "-g", "digital_sinnature")
			stdErr := string(session.Err.Contents())

			Eventually(session).Should(Exit(1))
			Expect(stdErr).To(MatchRegexp(`The provided key usage 'digital_sinnature' was not known. Please update this value and retry your request.`))
		})

		It("should handle CAs whose names have lots of special characters", func() {
			madDogCAId := "dan:test/ing?danothertbe$in&the[stuff]=that@shouldn!"

			By("setting a value with lots of special characters", func() {
				session := runCommand("ca-generate", "-n", madDogCAId, "--common-name", generateUniqueCredentialName())
				Eventually(session).Should(Exit(0))
			})

			By("retrieving the value that was set", func() {
				session := runCommand("ca-get", "-n", madDogCAId)
				Eventually(session).Should(Exit(0))
			})
		})
	})
})


// https://golang.org/pkg/crypto/x509/#Certificate
// prefix should be "Certificate" or "Ca"
func CertFromPem(input string, prefix string) *x509.Certificate {
	r, _ := regexp.Compile(`(?s)` + prefix + `:\s+(-----BEGIN CERTIFICATE-----.*?-----END CERTIFICATE-----)`)

	pemByteArrays := r.FindAllSubmatch([]byte(input), 5)[0]

	lastPem := pemByteArrays[len(pemByteArrays)-1]

	block, _ := pem.Decode(lastPem)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	return cert
}
