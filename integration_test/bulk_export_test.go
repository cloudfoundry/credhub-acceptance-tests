package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = It("should export secrets without metadata", func() {
	credentialPathWithoutMetadata := "/" + GenerateUniqueCredentialName()
	credentialNameWithoutMetadata := credentialPathWithoutMetadata + "/" + "secret-without-metadata"

	By("setting a secret without metadata", func() {
		session := RunCommand("set", "-n", credentialNameWithoutMetadata, "-t", "value", "-v", credentialValue)
		Eventually(session).Should(Exit(0))
	})

	By("exporting a secret without metadata", func() {
		session := RunCommand("export", "-p", credentialPathWithoutMetadata)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: ` + credentialNameWithoutMetadata))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: FAKE-CREDENTIAL-VALUE`))
		Expect(stdOut).ToNot(ContainSubstring(`metadata:`))
	})

	By("deleting the secrets", func() {
		session := RunCommand("delete", "-n", credentialNameWithoutMetadata)
		Eventually(session).Should(Exit(0))
	})
})
