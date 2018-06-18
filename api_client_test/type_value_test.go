package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"code.cloudfoundry.org/credhub-cli/credhub"
)

var _ = Describe("Value Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-value")
		cred := values.Value("some string value")
		cred2 := values.Value("another string value")

		By("setting the value for the first time returns same value")
		value, err := credhubClient.SetValue(name, cred, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred))

		By("setting the value again without overwrite returns same value if the parameters are the same")
		value, err = credhubClient.SetValue(name, cred2, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred))

		By("overwriting the value with set")
		value, err = credhubClient.SetValue(name, cred2, credhub.Overwrite)
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
