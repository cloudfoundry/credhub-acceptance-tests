package remote_backend_test

import (
	"encoding/json"
	"fmt"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

type Permission struct {
	Uuid       string   `json: uuid`
	Actor      string   `json: actor`
	Path       string   `json: path`
	Operations []string `json: operations`
}

var _ = Describe("Permissions", func() {

	Context("get permissions v2", func() {
		It("returns a NOT IMPLEMENTED error", func() {
			session := RunCommand("curl", "-p", "/api/v2/permissions?path=some-path&actor=some-actor")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("This resource has not been implemented for this backend."))
		})
	})

	Context("post permissions v2", func() {
		It("creates and returns the permission", func() {
			var actor = "some-actor"
			var path = "/some-path"
			var data = fmt.Sprintf(`{"path": "%s", "actor": "%s", "operations": ["read", "write"]}`, path, actor)
			session := RunCommand("curl", "-p", "/api/v2/permissions", "-X", "POST", "-d", data)
			Expect(session).Should(Exit(0))

			var permission Permission
			err := json.Unmarshal(session.Out.Contents(), &permission)
			Expect(err).NotTo(HaveOccurred())

			Expect(permission).Should(Equal(
				Permission{
					Uuid:       permission.Uuid,
					Actor:      actor,
					Path:       path,
					Operations: []string{"read", "write"},
				}))
		})
	})

	Context("put permissions v2", func() {
		It("updates (overwrite) the operations and returns the permission", func() {
			var actor = "some-actor"
			var path = GenerateUniqueCredentialName()
			var uuid = seedPermission(actor, path)

			var data = fmt.Sprintf(`{"path": "%s", "actor": "%s", "operations": ["read", "write"]}`, path, actor)
			session := RunCommand("curl", "-p", "/api/v2/permissions/"+uuid, "-X", "PUT", "-d", data)
			Expect(session).Should(Exit(0))

			var permission Permission
			err := json.Unmarshal(session.Out.Contents(), &permission)
			Expect(err).NotTo(HaveOccurred())

			Expect(permission).Should(Equal(
				Permission{
					Uuid:       uuid,
					Actor:      actor,
					Path:       "/" + path,
					Operations: []string{"read", "write"},
				}))
		})
	})

	Context("patch permissions v2", func() {
		It("updates the operations and returns the permission", func() {
			var actor = "some-actor"
			var path = GenerateUniqueCredentialName()
			var uuid = seedPermission(actor, path)

			session := RunCommand("curl", "-p", "/api/v2/permissions/"+uuid, "-X", "PATCH", "-d", `{"operations": ["write_acl"]}`)
			Expect(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring(`"uuid": `))

			var permission Permission
			err := json.Unmarshal(session.Out.Contents(), &permission)
			Expect(err).NotTo(HaveOccurred())

			Expect(permission).Should(Equal(
				Permission{
					Uuid:       uuid,
					Actor:      actor,
					Path:       "/" + path,
					Operations: []string{"write_acl"},
				}))
		})
	})
})

func seedPermission(actor, path string) string {
	var data = fmt.Sprintf(`{"path": "%s", "actor": "%s", "operations": ["read", "write", "read_acl"]}`, path, actor)
	session := RunCommand("curl", "-p", "/api/v2/permissions", "-X", "POST", "-d", data)

	var permission Permission
	err := json.Unmarshal(session.Out.Contents(), &permission)
	Expect(err).NotTo(HaveOccurred())

	return permission.Uuid
}
