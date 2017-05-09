package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Creating a User", func() {
	Describe("User Generation", func() {
		name := GenerateUniqueCredentialName()

		It("should create a user", func() {
			session := RunCommand("generate", "-n", name, "-t", "user")
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`name: /` + name))
			Expect(stdOut).To(ContainSubstring(`type: user`))
			Expect(stdOut).To(MatchRegexp(`username: \S*`))
			Expect(stdOut).To(MatchRegexp(`password: \S*\d`))
			Expect(stdOut).To(MatchRegexp(`password_hash: \$6\$.+\$.+`))
		})

		It("should get the created user", func() {
			session := RunCommand("get", "-n", name)
			stdOut := string(session.Out.Contents())

			Eventually(session).Should(Exit(0))

			Expect(stdOut).To(ContainSubstring(`name: /` + name))
			Expect(stdOut).To(ContainSubstring(`type: user`))
			Expect(stdOut).To(MatchRegexp(`username: \S*`))
			Expect(stdOut).To(MatchRegexp(`password: \S*\d`))
			Expect(stdOut).To(MatchRegexp(`password_hash: \$6\$.+\$.+`))
		})
	})
})
