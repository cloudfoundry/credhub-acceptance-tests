package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Password Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-password")
		generateParameters := generate.Password{Length: 10}

		By("generate a password with path " + name)
		password, err := credhubClient.GeneratePassword(name, generateParameters, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(generateParameters.Length))
		generatedPassword := password.Value

		By("generate the password again without overwrite returns same password")
		password, err = credhubClient.GeneratePassword(name, generateParameters, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(generatedPassword))

		By("overwriting the password with generate")
		password, err = credhubClient.GeneratePassword(name, generateParameters, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(generateParameters.Length))
		Expect(password.Value).ToNot(Equal(generatedPassword))

		newPassword := values.Password("some-password")

		By("setting the password again overwrites previous password")
		password, err = credhubClient.SetPassword(name, newPassword)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(newPassword))

		By("getting the password")
		password, err = credhubClient.GetLatestPassword(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(newPassword))

		By("deleting the password")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestPassword(name)
		Expect(err).To(HaveOccurred())
	})
})
