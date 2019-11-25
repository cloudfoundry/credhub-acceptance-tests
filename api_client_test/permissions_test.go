package acceptance_test

import (
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"code.cloudfoundry.org/credhub-cli/credhub/permissions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Getting Credentials", func() {
	It("Adds permission", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-password")

		_, err := credhubClient.GeneratePassword(name, generate.Password{}, credhub.Overwrite)
		Expect(err).NotTo(HaveOccurred())

		newPermission := permissions.Permission{
			Actor:      "some-actor",
			Operations: []string{"read"},
			Path:       name,
		}

		resp, err := credhubClient.AddPermission(newPermission.Path, newPermission.Actor, newPermission.Operations)
		Expect(resp.Actor).To(Equal("some-actor"))
		Expect(resp.Operations).To(Equal([]string{"read"}))
		Expect(resp.Path).To(Equal(name))
		Expect(err).NotTo(HaveOccurred())

		fetchedPermission, err := credhubClient.GetPermissionByUUID(resp.UUID)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedPermission.Actor).To(Equal(resp.Actor))
		Expect(fetchedPermission.Operations).To(Equal(resp.Operations))
		Expect(fetchedPermission.Path).To(Equal(resp.Path))
		Expect(fetchedPermission.UUID).To(Equal(resp.UUID))

		fetchedPermission, err = credhubClient.GetPermissionByPathActor(resp.Path, resp.Actor)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedPermission.Actor).To(Equal(resp.Actor))
		Expect(fetchedPermission.Operations).To(Equal(resp.Operations))
		Expect(fetchedPermission.Path).To(Equal(resp.Path))
		Expect(fetchedPermission.UUID).To(Equal(resp.UUID))
	})

	Context("when permission does not exist", func() {
		It("sets new permission", func() {
			name := testCredentialPath(time.Now().UnixNano(), "some-password")

			_, err := credhubClient.GeneratePassword(name, generate.Password{}, credhub.Overwrite)
			Expect(err).NotTo(HaveOccurred())

			_, err = credhubClient.GetPermissionByPathActor("/testpath", "test-actor")
			Expect(err).To(HaveOccurred())
		})

	})

})
