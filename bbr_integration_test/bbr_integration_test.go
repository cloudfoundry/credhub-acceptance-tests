package bbr_integration

import (
	"fmt"

	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Backup and Restore", func() {
	var credentialName string
	var bbrTestPath = "bbr_test"
	var bbrDirectory string

	BeforeEach(func() {
		credentialName = fmt.Sprintf("%s/%s", bbrTestPath, test_helpers.GenerateUniqueCredentialName())

		By("authenticating and targeting against credhub")
		session := RunCommand("credhub", "login", "--client-name", config.ClientName, "--client-secret", config.ClientSecret, "--server", config.ApiUrl, "--skip-tls-validation")
		Eventually(session).Should(Exit(0))

		CleanupCredhub(bbrTestPath)

		var err error
		bbrDirectory, err = ioutil.TempDir("", "bbr")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		CleanupCredhub(bbrTestPath)
		CleanupArtifacts()
		os.RemoveAll(bbrDirectory)
	})

	It("restores credentials from a backup", func() {
		By("creating test credentials")
		credhubCredentialNames := createCredhubCredentials(bbrTestPath)

		By("adding a test password that will be edited after backup")
		session := RunCommand("credhub", "set", "--name", credentialName, "--type", "password", "-w", "originalsecret")
		Eventually(session).Should(Exit(0))

		By("running bbr backup")
		session = RunCommand("bbr", "deployment", "--username", config.Bosh.Client, "--password", config.Bosh.ClientSecret, "--deployment", config.DeploymentName, "--target", config.Bosh.Environment, "--ca-cert", config.Bosh.CaCertPath, "backup", "--artifact-path", bbrDirectory)
		Eventually(session).Should(Exit(0))

		By("asserting that the backup archive exists and contains a database dump file")
		session = RunCommand("sh", "-c", fmt.Sprintf("tar tf %s/%s*/*credhubdb.tar ./credhubdb_dump", bbrDirectory, config.DeploymentName))
		Eventually(session).Should(Exit(0))

		By("editing the test password")
		session = RunCommand("credhub", "set", "--name", credentialName, "--type", "password", "-w", "updatedsecret")
		Eventually(session).Should(Exit(0))

		editSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(editSession).Should(Exit(0))
		Eventually(editSession.Out).Should(gbytes.Say("value: updatedsecret"))

		By("deleting all of the other test credentials")
		for _, credentialName := range credhubCredentialNames {
			session = RunCommand("credhub", "delete", "-n", credentialName)
			Eventually(session).Should(Exit(0))
		}

		By("running bbr restore")
		session = RunCommand("sh", "-c",
			fmt.Sprintf("bbr deployment --username %s --password %s --deployment %s --target %s --ca-cert %s restore --artifact-path %s/%s*",
				config.Bosh.Client, config.Bosh.ClientSecret, config.DeploymentName, config.Bosh.Environment, config.Bosh.CaCertPath, bbrDirectory, config.DeploymentName),
		)
		Eventually(session).Should(Exit(0))

		By("checking if the test password was restored")
		getSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(getSession).Should(Exit(0))
		Eventually(getSession.Out).Should(gbytes.Say("value: originalsecret"))

		By("checking if the other test credentials were restored")
		findSession := RunCommand("credhub", "find")
		Eventually(findSession).Should(Exit(0))
		for _, credentialName := range credhubCredentialNames {
			Expect(findSession.Out.Contents()).To(ContainSubstring(credentialName))
		}
	})
})

func createCredhubCredentials(credentialPrefix string) []string {
	passwordName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())
	certificateName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())
	sshName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())
	rsaName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())
	jsonName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())
	valueName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())
	userName := fmt.Sprintf("%s/%s", credentialPrefix, test_helpers.GenerateUniqueCredentialName())

	session := RunCommand("credhub", "generate", "--name", passwordName, "--type", "password")
	Eventually(session).Should(Exit(0))
	session = RunCommand("credhub", "generate", "--name", certificateName, "--type", "certificate", "-c", "cn", "--is-ca")
	Eventually(session).Should(Exit(0))
	session = RunCommand("credhub", "generate", "--name", sshName, "--type", "ssh")
	Eventually(session).Should(Exit(0))
	session = RunCommand("credhub", "generate", "--name", rsaName, "--type", "rsa")
	Eventually(session).Should(Exit(0))
	session = RunCommand("credhub", "generate", "--name", userName, "--type", "user")
	Eventually(session).Should(Exit(0))
	session = RunCommand("credhub", "set", "--name", jsonName, "--type", "json", "-v", `{"test": "secret"}`)
	Eventually(session).Should(Exit(0))
	session = RunCommand("credhub", "set", "--name", valueName, "--type", "value", "-v", "some-value")
	Eventually(session).Should(Exit(0))

	return []string{
		passwordName,
		certificateName,
		sshName,
		valueName,
		jsonName,
		rsaName,
		userName,
	}
}

func CleanupCredhub(path string) {
	By("Cleaning up credhub bbr test passwords")
	RunCommand(
		"sh", "-c",
		fmt.Sprintf("credhub find -p /%s --output-json | jq -r '.credentials[].name' | xargs -IN credhub delete --name N", path),
	)
}

func CleanupArtifacts() {
	By("Cleaning up bbr test artifacts")
	RunCommand("rm", "-rf", "credhubdb_dump")
	RunCommand("sh", "-c", fmt.Sprintf("rm -rf %s*Z", config.DirectorHost))
}
