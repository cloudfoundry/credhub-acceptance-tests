package certs

import (
	"crypto/x509"
	"fmt"

	"github.com/onsi/gomega/types"
)

func BeValidSelfSignedCert() types.GomegaMatcher {
	return &validSelfSignedCertMatcher{}
}

type validSelfSignedCertMatcher struct {
	validationError error
}

func (matcher *validSelfSignedCertMatcher) Match(actual interface{}) (bool, error) {
	cert, ok := actual.(*x509.Certificate)
	if !ok {
		return false, fmt.Errorf("BeValidSelfSignedCert matcher expects an x509.Certificate")
	}

	roots := x509.NewCertPool()
	roots.AddCert(cert)
	_, matcher.validationError = cert.Verify(x509.VerifyOptions{Roots: roots})
	return matcher.validationError == nil, nil
}

func (matcher *validSelfSignedCertMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected certificate validation to succeed, but got:\n\t%s", matcher.validationError.Error())
}

func (matcher *validSelfSignedCertMatcher) NegatedFailureMessage(actual interface{}) string {
	return "Expected certificate validation to fail, but it succeeded"
}

func BeValidCertSignedBy(expectedCA interface{}) types.GomegaMatcher {
	return &validCertSignedByMatcher{
		expectedCA: expectedCA,
	}
}

type validCertSignedByMatcher struct {
	expectedCA      interface{}
	validationError error
}

func (matcher *validCertSignedByMatcher) Match(actual interface{}) (bool, error) {
	cert, ok := actual.(*x509.Certificate)
	if !ok {
		return false, fmt.Errorf("BeValidCertSignedBy matcher expects an x509.Certificate")
	}

	expectedCABytes, ok := matcher.expectedCA.([]byte)
	if !ok {
		return false, fmt.Errorf("BeValidCertSignedBy matcher expects []byte of PEM-encoded CA")
	}

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(expectedCABytes)
	_, matcher.validationError = cert.Verify(x509.VerifyOptions{Roots: roots})
	return matcher.validationError == nil, nil
}

func (matcher *validCertSignedByMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected certificate validation to succeed, but got:\n\t%s", matcher.validationError.Error())
}

func (matcher *validCertSignedByMatcher) NegatedFailureMessage(actual interface{}) string {
	return "Expected certificate not to be valid and signed, but it was"
}

func FailCertValidationWithMessage(expectedMessage interface{}) types.GomegaMatcher {
	return &failCertValidationWithMessageMatcher{
		expectedMessage: expectedMessage,
	}
}

type failCertValidationWithMessageMatcher struct {
	expectedMessage interface{}
	validationError error
}

func (matcher *failCertValidationWithMessageMatcher) Match(actual interface{}) (bool, error) {
	cert, ok := actual.(*x509.Certificate)
	if !ok {
		return false, fmt.Errorf("BeValidSelfSignedCert matcher expects an x509.Certificate")
	}

	_, matcher.validationError = cert.Verify(x509.VerifyOptions{})
	return matcher.validationError != nil && matcher.validationError.Error() == matcher.expectedMessage, nil
}

func (matcher *failCertValidationWithMessageMatcher) FailureMessage(actual interface{}) string {
	if matcher.validationError == nil {
		return "Expected certificate validation to fail, but it succeeded"
	}
	return fmt.Sprintf("Expected certificate validation to fail with message:\n\t%s\nbut got:\n\t%s", matcher.expectedMessage, matcher.validationError.Error())
}

func (matcher *failCertValidationWithMessageMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected certificate validation not to fail with message:\n\t%s\nbut it did", matcher.expectedMessage)
}
