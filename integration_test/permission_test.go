package integration_test

import (
  "regexp"

  . "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/onsi/gomega/gexec"
)

var _ = Describe("Permission Test", func() {
  Context("Set Permmission is called on permission that does not exist", func(() {
    It("creates a new permmission", func() {
      session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write")
      Eventually(session).Should(Exit(0))
      Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "1234",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write"]
				}
				`))
    })
    })
  Context("Set Permission is called on permission that does exist", func() {
    It("updates existing permission", func() {
      session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
      Eventually(session).Should(Exit(0))
      Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "1234",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write", "delete"]
				}
				`))
    })
  })
})
