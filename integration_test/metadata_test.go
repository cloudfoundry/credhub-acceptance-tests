package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = It("should set a new secret with and without metadata", func() {
	supported, err := serverSupportsMetadata()
	Expect(err).NotTo(HaveOccurred())
	if !supported {
		Skip("Server does not support metadata")
	}

	credentialName1 := GenerateUniqueCredentialName() + "-with-metadata"
	credentialName2 := GenerateUniqueCredentialName() + "-without-metadata"

	By("setting a new secret with metadata", func() {
		session := RunCommand("set", "-n", credentialName1, "-t", "value", "-v", credentialValue, "--metadata", `{"some":"metadata"}`)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
  some: metadata
`))
	})

	By("setting a new secret without metadata", func() {
		session := RunCommand("set", "-n", credentialName2, "-t", "value", "-v", credentialValue)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).NotTo(ContainSubstring("metadata:"))
	})

	By("deleting the secrets", func() {
		session := RunCommand("delete", "-n", credentialName1)
		Eventually(session).Should(Exit(0))
		session = RunCommand("delete", "-n", credentialName2)
		Eventually(session).Should(Exit(0))
	})
})

func serverSupportsMetadata() (bool, error) {
	serverVersion, err := getServerVersion()
	if err != nil {
		return false, err
	}
	checkVersion, err := version.NewVersion("2.6.0")
	if err != nil {
		return false, err
	}
	serverVersion.GreaterThan(checkVersion)

	return serverVersion.GreaterThan(checkVersion) || serverVersion.Equal(checkVersion), nil
}
