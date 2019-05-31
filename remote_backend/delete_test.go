package remote_backend_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Delete", func() {
	Describe("get by name", func() {
		It("deletes a cred by name", func() {
			name := "/some-value"
			value := "some-random-value"

			session := RunCommand("set", "-t", "value", "-n", name, "-v", value)
			Expect(session).Should(Exit(0))

			session = RunCommand("delete", "-n", name)
			Expect(session).Should(Exit(0))

			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("Credential successfully deleted"))
		})
	})
})
