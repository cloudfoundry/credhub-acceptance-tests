package acceptance_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

var _ = Describe("User Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-user")
		generateParameters := generate.User{Length: 10}

		By("generate a user with path " + name)
		user, err := credhubClient.GenerateUser(name, generateParameters, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.Password).To(HaveLen(generateParameters.Length))
		generatedUser := user.Value

		By("generate the user again without overwrite returns same user")
		user, err = credhubClient.GenerateUser(name, generateParameters, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value).To(Equal(generatedUser))

		By("overwriting the user with generate")
		user, err = credhubClient.GenerateUser(name, generateParameters, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.Password).To(HaveLen(generateParameters.Length))
		Expect(user.Value).ToNot(Equal(generatedUser))

		username := "name"
		newUser := values.User{Username: username, Password: "password"}

		By("setting the user again overwrites previous user")
		user, err = credhubClient.SetUser(name, newUser)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.User).To(Equal(newUser))

		By("getting the user")
		user, err = credhubClient.GetLatestUser(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.User).To(Equal(newUser))

		By("deleting the user")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestUser(name)
		Expect(err).To(HaveOccurred())
	})
})
