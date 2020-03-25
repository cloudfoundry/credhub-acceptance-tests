package acceptance_test

import (
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata", func() {
	BeforeEach(func() {
		supported, err := serverSupportsMetadata()
		Expect(err).NotTo(HaveOccurred())
		if !supported {
			Skip("Server does not support metadata")
		}
	})

	It("should set a new secret with and without metadata", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-value")
		cred := values.Value("some string value")
		cred2 := values.Value("another string value")

		By("setting the value for the first time returns new value")
		value, err := credhubClient.SetValue(name, cred)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred))

		By("setting the value again overwrites previous value")
		value, err = credhubClient.SetValue(name, cred2)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred2))

		By("getting the value")
		value, err = credhubClient.GetLatestValue(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred2))

		By("deleting the value")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestValue(name)
		Expect(err).To(HaveOccurred())
	})
})

func serverSupportsMetadata() (bool, error) {
	serverVersion, err := getServerVersion()
	if err != nil {
		return false, err
	}

	checkVersion, err := version.NewVersion("2.6.0")
	if err != nil {
		return false, err
	}
	serverVersion.GreaterThan(checkVersion)

	return serverVersion.GreaterThan(checkVersion) || serverVersion.Equal(checkVersion), nil
}
