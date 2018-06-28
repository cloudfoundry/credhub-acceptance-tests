package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"code.cloudfoundry.org/credhub-cli/credhub"
)

var _ = Describe("Password Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-password")
		generatePassword := generate.Password{Length: 10}

		By("generate a password with path " + name)
		password, err := credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
		firstGeneratedPassword := password.Value

		By("generate the password again without overwrite returns same password")
		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(firstGeneratedPassword))

		By("setting the password again without overwrite returns same password")
		password, err = credhubClient.SetPassword(name, values.Password("some-password"), credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(firstGeneratedPassword))

		By("overwriting the password with generate")
		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
		Expect(password.Value).ToNot(Equal(firstGeneratedPassword))

		By("overwriting the password with set")
		password, err = credhubClient.SetPassword(name, values.Password("some-password"), credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(values.Password("some-password")))

		By("getting the password")
		password, err = credhubClient.GetLatestPassword(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(values.Password("some-password")))

		By("deleting the password")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestPassword(name)
		Expect(err).To(HaveOccurred())
	})
})