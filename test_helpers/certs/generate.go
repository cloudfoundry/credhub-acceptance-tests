package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

const RsaKeySize = 4096

type CertOptions struct {
	CommonName         string
	OrganizationalUnit string
	IsCA               bool
	NotBefore          time.Time
	NotAfter           time.Time
}

func GenerateSigned(certOptions CertOptions, caCert []byte, caKey []byte) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, RsaKeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key: %s", err)
	}
	csrDerBytes, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{}, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CSR: %s", err.Error())
	}
	csr, err := x509.ParseCertificateRequest(csrDerBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CSR: %s", err.Error())
	}

	template, err := generateCertificateTemplate(certOptions)
	if err != nil {
		return nil, nil, err
	}

	caTLS, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load CA key pair: %s", err)
	}
	ca, err := x509.ParseCertificate(caTLS.Certificate[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %s", err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, ca, csr.PublicKey, caTLS.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %s", err)
	}

	return encodeCertAndKey(derBytes, key)
}

func GenerateSelfSigned(certOptions CertOptions) ([]byte, []byte, error) {
	template, err := generateCertificateTemplate(certOptions)
	if err != nil {
		return nil, nil, err
	}

	key, err := rsa.GenerateKey(rand.Reader, RsaKeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key: %s", err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %s", err)
	}

	return encodeCertAndKey(derBytes, key)
}

func generateCertificateTemplate(certOptions CertOptions) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	notBefore, notAfter, err := calculateExpiryDates(certOptions)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Subject:      pkix.Name{CommonName: certOptions.CommonName},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
	}
	if certOptions.OrganizationalUnit != "" {
		template.Subject.OrganizationalUnit = []string{certOptions.OrganizationalUnit}
	}
	if certOptions.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		template.BasicConstraintsValid = true
	}

	return template, nil
}

func calculateExpiryDates(certOptions CertOptions) (time.Time, time.Time, error) {
	notBefore := time.Now()
	if !certOptions.NotBefore.IsZero() {
		notBefore = certOptions.NotBefore
	}
	notAfter := notBefore.Add(time.Hour * 24 * 30)
	if !certOptions.NotAfter.IsZero() {
		notAfter = certOptions.NotAfter
	}
	if notBefore.After(notAfter) {
		return time.Time{}, time.Time{}, fmt.Errorf("NotBefore (%s) must be earlier than NotAfter (%s)", notBefore, notAfter)
	}
	return notBefore, notAfter, nil
}

func encodeCertAndKey(derBytes []byte, key *rsa.PrivateKey) ([]byte, []byte, error) {
	var certPem, keyPem bytes.Buffer
	if err := pem.Encode(&certPem, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, fmt.Errorf("failed to encode certificate: %s", err)
	}
	if err := pem.Encode(&keyPem, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return nil, nil, fmt.Errorf("failed to encode key: %s", err)
	}

	return certPem.Bytes(), keyPem.Bytes(), nil
}
