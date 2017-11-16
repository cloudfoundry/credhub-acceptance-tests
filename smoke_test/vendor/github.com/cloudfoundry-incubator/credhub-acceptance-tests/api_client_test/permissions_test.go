package acceptance_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Getting Credentials", func() {
	Specify("GetLatestVersion", func() {
		name := testCredentialPath("some-password")

		_, err := credhubClient.GeneratePassword(name, generate.Password{}, credhub.Overwrite)
		Expect(err).NotTo(HaveOccurred())

		newPermission := permissions.Permission{
			Actor:      "some-actor",
			Operations: []string{"read"},
		}

		_, err = credhubClient.AddPermissions(name, []permissions.Permission{newPermission})
		Expect(err).NotTo(HaveOccurred())

		fetchedPermissions, err := credhubClient.GetPermissions(name)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedPermissions).To(ContainElement(newPermission))
	})
})
