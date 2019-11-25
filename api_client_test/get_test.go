package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
)

var _ = Describe("Getting Credentials", func() {
	Specify("GetLatestVersion", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		firstPassword := string(password.Value)

		credhubClient.GetLatestPassword(name)

		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		secondPassword := string(password.Value)

		credential, err := credhubClient.GetLatestVersion(name)

		Expect(credential.Name).To(Equal(name))
		Expect(credential.Value).To(Equal(secondPassword))
		Expect(credential.Value).ToNot(Equal(firstPassword))

		Expect(credhubClient.Delete(name)).NotTo(HaveOccurred())
	})

	Specify("GetNVersions", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())

		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		secondPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		thirdPassword := string(password.Value)

		credentials, err := credhubClient.GetNVersions(name, 2)

		Expect(credentials[0].Name).To(Equal(name))
		Expect(credentials[0].Value).To(Equal(thirdPassword))
		Expect(credentials[1].Value).To(Equal(secondPassword))

		Expect(credhubClient.Delete(name)).NotTo(HaveOccurred())
	})

	Specify("GetAllVersions", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		firstPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		secondPassword := string(password.Value)

		password, err = credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		thirdPassword := string(password.Value)

		credentials, err := credhubClient.GetAllVersions(name)

		Expect(credentials[0].Name).To(Equal(name))
		Expect(credentials[0].Value).To(Equal(thirdPassword))
		Expect(credentials[1].Value).To(Equal(secondPassword))
		Expect(credentials[2].Value).To(Equal(firstPassword))

		Expect(credhubClient.Delete(name)).NotTo(HaveOccurred())
	})

	Specify("GetById", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-password")

		generatePassword := generate.Password{}

		By("generate a password with path " + name)

		password, err := credhubClient.GeneratePassword(name, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).ToNot(BeEmpty())
		firstPasswordId := string(password.Id)

		credential, err := credhubClient.GetById(firstPasswordId)

		Expect(credential.Name).To(Equal(name))
		Expect(credential.Value).To(Equal(string(password.Value)))

		Expect(credhubClient.Delete(name)).NotTo(HaveOccurred())
	})
})
