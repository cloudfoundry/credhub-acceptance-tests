package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/utilities"
	"time"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

type Permission struct {
	Actor      string   `json: actor`
	Path       string   `json: path`
	Operations []string `json: operations`
}

var _ = Describe("Permission Test", func() {
	var actor string
	var path string

	BeforeEach(func() {
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
		actor = "actor-" + timestamp
		path = "/path-" + timestamp
	})

	Context("Set Permission", func() {
		Context("Set Permission is called on permission that does not exist", func() {
			Context("when output json flag is used", func() {
				It("creates a new permission", func() {
					session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write", "-j")
					Eventually(session).Should(Exit(0))
					var permission Permission
					err := json.Unmarshal(session.Out.Contents(), &permission)
					Expect(err).NotTo(HaveOccurred())

					Expect(permission).Should(Equal(
						Permission{
							Actor:      actor,
							Path:       path,
							Operations: []string{"read", "write"},
						}))
				})
			})
			It("creates a new permission", func() {
				session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				err := utilities.GenerateAutoDoc(session)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(Exit(0))

				Eventually(string(session.Out.Contents())).Should(ContainSubstring("uuid: "))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("actor: " + actor))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("path: " + path))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring(`operations:
- read
- write`))
			})
		})
		Context("Set Permission is called on permission that exists", func() {
			Context("when output json flag is used", func() {
				It("updates existing permission", func() {
					session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
					Eventually(session).Should(Exit(0))
					session = RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write, delete", "-j")
					Eventually(session).Should(Exit(0))

					var permission Permission
					err := json.Unmarshal(session.Out.Contents(), &permission)
					Expect(err).NotTo(HaveOccurred())

					Expect(permission).Should(Equal(
						Permission{
							Actor:      actor,
							Path:       path,
							Operations: []string{"read", "write", "delete"},
						}))
				})
			})
			It("updates existing permission", func() {
				session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				Eventually(session).Should(Exit(0))
				session = RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write, delete")
				Eventually(session).Should(Exit(0))

				Eventually(string(session.Out.Contents())).Should(ContainSubstring("uuid: "))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("actor: " + actor))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("path: " + path))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring(`operations:
- read
- write`))
			})
		})
	})
	Context("Get Permission", func() {
		Context("Get permission called on permission that does not exist", func() {
			It("throws an error", func() {
				session := RunCommand("get-permission", "-a", actor, "-p", path)
				Eventually(session).Should(Exit(1))
			})
		})
		Context("Get permission called on permission that exists", func() {
			Context("when output json flag is used", func() {
				It("returns the permission", func() {
					session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
					Eventually(session).Should(Exit(0))
					session = RunCommand("get-permission", "-a", actor, "-p", path, "-j")
					Eventually(session).Should(Exit(0))

					var permission Permission
					err := json.Unmarshal(session.Out.Contents(), &permission)
					Expect(err).NotTo(HaveOccurred())

					Expect(permission).Should(Equal(
						Permission{
							Actor:      actor,
							Path:       path,
							Operations: []string{"read", "write"},
						}))
				})
			})
			It("returns the permission", func() {
				session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				Eventually(session).Should(Exit(0))
				session = RunCommand("get-permission", "-a", actor, "-p", path)
				err := utilities.GenerateAutoDoc(session)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(Exit(0))

				Eventually(string(session.Out.Contents())).Should(ContainSubstring("uuid: "))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("actor: " + actor))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("path: " + path))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring(`operations:
- read
- write`))
			})
		})
	})
	Context("Delete Permission", func() {
		Context("Delete permission called on permission that does not exist", func() {
			It("throws an error", func() {
				session := RunCommand("delete-permission", "-a", actor, "-p", path)
				Eventually(session).Should(Exit(1))
			})
		})
		Context("Delete permission called on permission that exists", func() {
			It("returns the deleted permission", func() {
				session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				Eventually(session).Should(Exit(0))
				session = RunCommand("delete-permission", "-a", actor, "-p", path)
				err := utilities.GenerateAutoDoc(session)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(Exit(0))

				Eventually(string(session.Out.Contents())).Should(ContainSubstring("uuid: "))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("actor: " + actor))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring("path: " + path))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring(`operations:
- read
- write`))
			})
		})
	})
})
