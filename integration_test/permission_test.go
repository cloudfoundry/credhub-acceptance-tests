package integration_test

import (
	"encoding/json"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)


type Permission struct {
	Actor string `json: actor`
	Path string		`json: path`
	Operations []string `json: operations`
}

var _ = Describe("Permission Test", func() {
	Context("Set Permission is called on permission that does not exist", func() {
		It("creates a new permission", func() {
			session := RunCommand("set-permission", "-a", "some-actor", "-p", "/some-path", "-o", "read, write")
			var permission Permission
			_ = json.Unmarshal(session.Out.Contents(), &permission)
			Expect(permission).Should(Equal(
				Permission{
					Actor: "some-actor",
					Path: "/some-path",
					Operations: []string{"read", "write"},
				}))
			Eventually(session).Should(Exit(0))
		})
	})
	Context("Set Permission is called on permission that does exist", func() {
		It("updates existing permission", func() {
			RunCommand("set-permission", "-a", "some-actor", "-p", "/some-path", "-o", "read, write")
			session := RunCommand("set-permission", "-a", "some-actor", "-p", "/some-path", "-o", "read, write, delete")
			var permission Permission
			_ = json.Unmarshal(session.Out.Contents(), &permission)
			Expect(permission).Should(Equal(
				Permission{
					Actor: "some-actor",
					Path: "/some-path",
					Operations: []string{"read", "write", "delete"},
				}))
			Eventually(session).Should(Exit(0))
		})
	})
})
