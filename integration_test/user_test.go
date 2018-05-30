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

		Describe("With default parameters", func() {
			It("should generate a user", func() {
				RunCommand("generate", "-n", name, "-t", "user")
				session := RunCommand("get", "-n", name)
				stdOut := string(session.Out.Contents())

				By("generating the credential first", func() {
					Eventually(session).Should(Exit(0))

					Expect(stdOut).To(ContainSubstring(`name: /` + name))
					Expect(stdOut).To(ContainSubstring(`type: user`))
					Expect(stdOut).To(MatchRegexp(`username: \S*`))
					Expect(stdOut).To(MatchRegexp(`password: \S*\d`))
					Expect(stdOut).To(MatchRegexp(`password_hash: \$6\$.+\$.+`))
				})

				By("getting the generated credential", func() {
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

		Describe("with parameters", func() {
			It("should generate a user with password of length 50", func() {
				RunCommand("generate", "-n", name, "-t", "user", "--length", "50")
				session := RunCommand("get", "-n", name)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`name: /` + name))
				Expect(stdOut).To(ContainSubstring(`type: user`))
				Expect(stdOut).To(MatchRegexp(`username: \S*`))
				Expect(stdOut).To(MatchRegexp(`password: \S{50}\b`))
				Expect(stdOut).To(MatchRegexp(`password_hash: \$6\$.+\$.+`))
			})
		})

		Describe("with provided username", func() {
			It("should generate a password, but not the username", func() {
				username := "test-username"
				RunCommand("generate", "-n", name, "-t", "user", "--username", username)
				session := RunCommand("get", "-n", name)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`name: /` + name))
				Expect(stdOut).To(ContainSubstring(`type: user`))
				Expect(stdOut).To(ContainSubstring(`username: ` + username))
				Expect(stdOut).To(MatchRegexp(`password: \S*\d`))
				Expect(stdOut).To(MatchRegexp(`password_hash: \$6\$.+\$.+`))
			})
		})
	})

	Describe("Setting a User value", func() {
		name := GenerateUniqueCredentialName()

		Describe("including all parameters", func() {
			It("should set the user value", func() {
				username := "test"
				password := "password"
				RunCommand("set", "-n", name, "-t", "user", "-z", username, "-w", password)
				session := RunCommand("get", "-n", name)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))

				Expect(stdOut).To(ContainSubstring(`name: /` + name))
				Expect(stdOut).To(ContainSubstring(`type: user`))
				Expect(stdOut).To(ContainSubstring(`username: ` + username))
				Expect(stdOut).To(ContainSubstring(`password: ` + password))
				Expect(stdOut).To(MatchRegexp(`password_hash: \$6\$.+\$.+`))
			})
		})
	})
})
