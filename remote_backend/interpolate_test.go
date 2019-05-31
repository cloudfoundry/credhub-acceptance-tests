package remote_backend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Interpolate", func() {

	It("returns a NOT IMPLEMENTED error", func() {
		vcapServicesString := `{"someOuterKey":[{"credentials": {"credhub-ref":"((/some-cred))"},"anotherKey":"some other value that does not change"}]}`
		session := RunCommand("curl", "-XPOST", "-p", "/api/v1/interpolate", "-d", vcapServicesString)
		Expect(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring("This resource has not been implemented for this backend."))
	})

})
