package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Import test", func() {
	It("should import credentials from a file", func() {
		session := RunCommand("import", "-f", "../test_helpers/bulk_import.yml")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - agent`))
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - director`))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))
	})

	It("should save the credentials on CredHub", func() {
		session := RunCommand("get", "-n", "/director/deployment/blobstore - agent" )
		Eventually(session).Should(Exit(0))
		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - agent`))
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))

		session = RunCommand("get", "-n", "/director/deployment/blobstore - director" )
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - director`))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))
	})

})
