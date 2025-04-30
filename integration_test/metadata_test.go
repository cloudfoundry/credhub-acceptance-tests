package integration_test

import (
	"encoding/json"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	"github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = It("should generate a new secret with and without metadata", func() {
	supported, err := serverSupportsMetadata()
	Expect(err).NotTo(HaveOccurred())
	if !supported {
		Skip("Server does not support metadata")
	}

	credentialName1 := GenerateUniqueCredentialName() + "-with-metadata"
	credentialName2 := GenerateUniqueCredentialName() + "-without-metadata"

	By("generating a new secret with metadata", func() {
		session := RunCommand("generate", "-n", credentialName1, "-t", "password", "--metadata", `{"some":"metadata"}`)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
    some: metadata
`))
	})

	By("setting a new secret without metadata", func() {
		session := RunCommand("generate", "-n", credentialName2, "-t", "password")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).NotTo(ContainSubstring("metadata:"))
	})

	By("getting a secret with metadata", func() {
		session := RunCommand("get", "-n", credentialName1)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
    some: metadata
`))
	})

	By("getting a secret with metadata with --output-json flag", func() {
		session := RunCommand("get", "-n", credentialName1, "--output-json")
		Eventually(session).Should(Exit(0))

		var output map[string]interface{}
		err = json.Unmarshal(session.Out.Contents(), &output)
		Expect(err).NotTo(HaveOccurred())

		Expect(output).To(HaveKeyWithValue("name", "/"+credentialName1))
		Expect(output).To(HaveKeyWithValue("type", "password"))
		Expect(output["value"]).ToNot(BeEmpty())
		Expect(output).To(HaveKeyWithValue("metadata", map[string]interface{}{"some": "metadata"}))
	})

	By("getting a secret without metadata", func() {
		session := RunCommand("get", "-n", credentialName2)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).ToNot(ContainSubstring(`metadata:`))
	})

	By("getting a secret without metadata with --output-json flag", func() {
		session := RunCommand("get", "-n", credentialName2, "--output-json")
		Eventually(session).Should(Exit(0))

		var output map[string]interface{}
		err = json.Unmarshal(session.Out.Contents(), &output)
		Expect(err).NotTo(HaveOccurred())

		Expect(output).To(HaveKeyWithValue("name", "/"+credentialName2))
		Expect(output).To(HaveKeyWithValue("type", "password"))
		Expect(output["value"]).ToNot(BeEmpty())
		Expect(output).To(HaveKeyWithValue("metadata", BeNil()))
	})

	By("deleting the secrets", func() {
		session := RunCommand("delete", "-n", credentialName1)
		Eventually(session).Should(Exit(0))
		session = RunCommand("delete", "-n", credentialName2)
		Eventually(session).Should(Exit(0))
	})
})

var _ = It("should regenerate a secret with and without metadata", func() {
	supported, err := serverSupportsMetadata()
	Expect(err).NotTo(HaveOccurred())
	if !supported {
		Skip("Server does not support metadata")
	}

	credentialName1 := GenerateUniqueCredentialName() + "-with-metadata"
	credentialName2 := GenerateUniqueCredentialName() + "-without-metadata"

	session := RunCommand("generate", "-n", credentialName1, "-t", "password", "--metadata", `{"some":"metadata"}`)
	Eventually(session).Should(Exit(0))

	session2 := RunCommand("generate", "-n", credentialName2, "-t", "password", "--metadata", `{"some":"other-metadata"}`)
	Eventually(session2).Should(Exit(0))

	By("regenerating a secret with metadata", func() {
		session := RunCommand("regenerate", "-n", credentialName1, "--metadata", `{"some":"regenerated-metadata"}`)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
    some: regenerated-metadata
`))
	})

	By("getting a regenerated secret with metadata", func() {
		session := RunCommand("get", "-n", credentialName1)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
    some: regenerated-metadata
`))
	})

	By("getting a regenerated secret with metadata with --output-json flag", func() {
		session := RunCommand("get", "-n", credentialName1, "--output-json")
		Eventually(session).Should(Exit(0))

		var output map[string]interface{}
		err = json.Unmarshal(session.Out.Contents(), &output)
		Expect(err).NotTo(HaveOccurred())

		Expect(output).To(HaveKeyWithValue("name", "/"+credentialName1))
		Expect(output).To(HaveKeyWithValue("type", "password"))
		Expect(output["value"]).ToNot(BeEmpty())
		Expect(output).To(HaveKeyWithValue("metadata", map[string]interface{}{"some": "regenerated-metadata"}))
	})

	By("regenerating a secret without metadata", func() {
		session := RunCommand("regenerate", "-n", credentialName2)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
  some: other-metadata
`))
	})

	By("getting a regenerated secret without metadata", func() {
		session := RunCommand("get", "-n", credentialName2)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
  some: other-metadata
`))
	})

	By("getting a regenerated secret without metadata with --output-json flag", func() {
		session := RunCommand("get", "-n", credentialName2, "--output-json")
		Eventually(session).Should(Exit(0))

		var output map[string]interface{}
		err = json.Unmarshal(session.Out.Contents(), &output)
		Expect(err).NotTo(HaveOccurred())

		Expect(output).To(HaveKeyWithValue("name", "/"+credentialName2))
		Expect(output).To(HaveKeyWithValue("type", "password"))
		Expect(output["value"]).ToNot(BeEmpty())
		Expect(output).To(HaveKeyWithValue("metadata", map[string]interface{}{"some": "other-metadata"}))
	})

	By("deleting the secrets", func() {
		session := RunCommand("delete", "-n", credentialName1)
		Eventually(session).Should(Exit(0))
		session = RunCommand("delete", "-n", credentialName2)
		Eventually(session).Should(Exit(0))
	})
})

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

	By("getting a secret with metadata", func() {
		session := RunCommand("get", "-n", credentialName1)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: FAKE-CREDENTIAL-VALUE`))
		Expect(stdOut).To(ContainSubstring(`
metadata:
    some: metadata
`))
	})

	By("getting a secret with metadata with --output-json flag", func() {
		session := RunCommand("get", "-n", credentialName1, "--output-json")
		Eventually(session).Should(Exit(0))

		var output map[string]interface{}
		err = json.Unmarshal(session.Out.Contents(), &output)
		Expect(err).NotTo(HaveOccurred())

		Expect(output).To(HaveKeyWithValue("name", "/"+credentialName1))
		Expect(output).To(HaveKeyWithValue("type", "value"))
		Expect(output).To(HaveKeyWithValue("value", "FAKE-CREDENTIAL-VALUE"))
		Expect(output).To(HaveKeyWithValue("metadata", map[string]interface{}{"some": "metadata"}))
	})

	By("getting a secret without metadata", func() {
		session := RunCommand("get", "-n", credentialName2)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: FAKE-CREDENTIAL-VALUE`))
		Expect(stdOut).ToNot(ContainSubstring(`metadata:`))
	})

	By("getting a secret without metadata with --output-json flag", func() {
		session := RunCommand("get", "-n", credentialName2, "--output-json")
		Eventually(session).Should(Exit(0))

		var output map[string]interface{}
		err = json.Unmarshal(session.Out.Contents(), &output)
		Expect(err).NotTo(HaveOccurred())

		Expect(output).To(HaveKeyWithValue("name", "/"+credentialName2))
		Expect(output).To(HaveKeyWithValue("type", "value"))
		Expect(output).To(HaveKeyWithValue("value", "FAKE-CREDENTIAL-VALUE"))
		Expect(output).To(HaveKeyWithValue("metadata", BeNil()))
	})

	By("deleting the secrets", func() {
		session := RunCommand("delete", "-n", credentialName1)
		Eventually(session).Should(Exit(0))
		session = RunCommand("delete", "-n", credentialName2)
		Eventually(session).Should(Exit(0))
	})
})

var _ = It("should export secrets with and without metadata", func() {
	supported, err := serverSupportsMetadata()
	Expect(err).NotTo(HaveOccurred())
	if !supported {
		Skip("Server does not support metadata")
	}

	credentialPathWithMetadata := "/" + GenerateUniqueCredentialName()
	credentialNameWithMetadata := credentialPathWithMetadata + "/" + "secret-with-metadata"

	credentialPathWithoutMetadata := "/" + GenerateUniqueCredentialName()
	credentialNameWithoutMetadata := credentialPathWithoutMetadata + "/" + "secret-without-metadata"

	By("setting a secret with metadata & a secret without metadata", func() {
		session := RunCommand("set", "-n", credentialNameWithMetadata, "-t", "value", "-v", credentialValue, "--metadata", `{"some":"metadata"}`)
		Eventually(session).Should(Exit(0))
		session = RunCommand("set", "-n", credentialNameWithoutMetadata, "-t", "value", "-v", credentialValue)
		Eventually(session).Should(Exit(0))
	})

	By("exporting a secret with metadata", func() {
		session := RunCommand("export", "-p", credentialPathWithMetadata)
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: ` + credentialNameWithMetadata))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: FAKE-CREDENTIAL-VALUE`))
		Expect(stdOut).To(ContainSubstring(`some: metadata`))
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
		session := RunCommand("delete", "-n", credentialNameWithMetadata)
		Eventually(session).Should(Exit(0))
		session = RunCommand("delete", "-n", credentialNameWithoutMetadata)
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
