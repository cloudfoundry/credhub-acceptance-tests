package remote_backend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"

)

var _ = Describe("Permissions", func() {

	Context("permissions v2", func(){
		It("returns a NOT IMPLEMENTED error", func() {
			session := RunCommand("curl", "-p", "/api/v2/permissions?path=some-path&actor=some-actor")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("This resource has not been implemented for this backend."))
		})
	})

	Context("permissions v1", func(){
		It("returns a NOT IMPLEMENTED error", func() {
			session := RunCommand("curl", "-p", "/api/v1/permissions?credential_name=some-cred")
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("This resource has not been implemented for this backend."))
		})

	})




})
