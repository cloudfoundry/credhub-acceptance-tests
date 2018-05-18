package bbr_integration

import (
	"fmt"

	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
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

	It("Successfully backs up and restores a Credhub release", func() {
		By("adding a test credential")
		session := RunCommand("credhub", "set", "--name", credentialName, "--type", "password", "-w", "originalsecret")
		Eventually(session).Should(Exit(0))

		By("running bbr backup")
		session = RunCommand("bbr", "deployment", "--username", config.Bosh.Client, "--password", config.Bosh.ClientSecret, "--deployment", config.DeploymentName, "--target", config.Bosh.Environment, "--ca-cert", config.Bosh.CaCertPath, "backup", "--artifact-path", bbrDirectory)
		Eventually(session).Should(Exit(0))

		By("asserting that the backup archive exists and contains a database dump file")
		session = RunCommand("sh", "-c", fmt.Sprintf("tar tf %s/%s*/*credhubdb.tar ./credhubdb_dump", bbrDirectory, config.DeploymentName))
		Eventually(session).Should(Exit(0))

		By("editing the test credential")
		session = RunCommand("credhub", "set", "--name", credentialName, "--type", "password", "-w", "updatedsecret")
		Eventually(session).Should(Exit(0))

		editSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(editSession).Should(Exit(0))
		Eventually(editSession.Out).Should(gbytes.Say("value: updatedsecret"))

		By("running bbr restore")
		session = RunCommand("sh", "-c",
			fmt.Sprintf("bbr deployment --username %s --password %s --deployment %s --target %s --ca-cert %s restore --artifact-path %s/%s*",
				config.Bosh.Client, config.Bosh.ClientSecret, config.DeploymentName, config.Bosh.Environment, config.Bosh.CaCertPath, bbrDirectory, config.DeploymentName),
		)
		Eventually(session).Should(Exit(0))

		By("checking if the test credentials was restored")
		getSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(getSession).Should(Exit(0))
		Eventually(getSession.Out).Should(gbytes.Say("value: originalsecret"))
	})
})

func CleanupCredhub(path string) {
	By("Cleaning up credhub bbr test passwords")
	RunCommand(
		"sh", "-c",
		fmt.Sprintf("credhub find -p /%s | tail -n +2 | cut -d\" \" -f1 | xargs -IN credhub delete --name N", path),
	)
}

func CleanupArtifacts() {
	By("Cleaning up bbr test artifacts")
	RunCommand("rm", "-rf", "credhubdb_dump")
	RunCommand("sh", "-c", fmt.Sprintf("rm -rf %s*Z", config.DirectorHost))
}
