package remote_backend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Regenerate", func() {
	Context("bulk regenerate", func() {
		It("returns a NOT IMPLEMENTED error", func() {
			session := RunCommand("curl", "-XPOST", "-p", "/api/v1/bulk-regenerate", "-d", `{"signed_by":"some-ca"}`)
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("This resource has not been implemented for this backend."))
		})
	})
	Context("regenerate", func() {
		It("returns a NOT IMPLEMENTED error", func() {
			session := RunCommand("curl", "-XPOST", "-p", "/api/v1/regenerate", "-d", `{"name":"some-cert"}`)
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("This resource has not been implemented for this backend."))
		})
	})
})
