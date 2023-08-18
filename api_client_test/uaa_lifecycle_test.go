package acceptance_test

import (
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Lifecycle", func() {
	Specify("lifecycle", func() {

		oauth := credhubClient.Auth.(*auth.OAuthStrategy)

		Expect(oauth.Login()).ToNot(HaveOccurred())
		Expect(oauth.Logout()).ToNot(HaveOccurred())

		Expect(oauth.Login()).ToNot(HaveOccurred())
		Expect(oauth.Logout()).ToNot(HaveOccurred())
	})
})
