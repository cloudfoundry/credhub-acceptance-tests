package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

var _ = Describe("Getting Credentials", func() {
	Specify("GetLatestVersion", func() {
		name := testCredentialPath("some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		firstPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		secondPassword := string(password.Value)

		credential, err := credhubClient.GetLatestVersion(name)

		Expect(credential.Name).To(Equal(name))
		Expect(credential.Value).To(Equal(secondPassword))
		Expect(credential.Value).ToNot(Equal(firstPassword))
	})

	Specify("GetNVersions", func() {
		name := testCredentialPath("some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())

		password, err = credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		secondPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		thirdPassword := string(password.Value)

		credentials, err := credhubClient.GetNVersions(name, 2)

		Expect(credentials[0].Name).To(Equal(name))
		Expect(credentials[0].Value).To(Equal(thirdPassword))
		Expect(credentials[1].Value).To(Equal(secondPassword))
	})

	Specify("GetAllVersions", func() {
		name := testCredentialPath("some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		firstPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		secondPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		thirdPassword := string(password.Value)

		credentials, err := credhubClient.GetAllVersions(name)

		Expect(credentials[0].Name).To(Equal(name))
		Expect(credentials[0].Value).To(Equal(thirdPassword))
		Expect(credentials[1].Value).To(Equal(secondPassword))
		Expect(credentials[2].Value).To(Equal(firstPassword))
	})

	Specify("GetById", func() {
		name := testCredentialPath("some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		firstPasswordId := string(password.Id)

		credential, err := credhubClient.GetById(firstPasswordId)

		Expect(credential.Name).To(Equal(name))
		Expect(credential.Value).To(Equal(string(password.Value)))
	})
})
