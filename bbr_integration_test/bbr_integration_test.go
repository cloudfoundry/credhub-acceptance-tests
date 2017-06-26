package bbr_integration

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Backup and Restore", func() {
	var credentialName string
	var bbrTestPath = "bbr_test"

	BeforeEach(func() {
		credentialName = fmt.Sprintf("%s/%s", bbrTestPath, test_helpers.GenerateUniqueCredentialName())

		By("authenticating against credhub")
		RunCommand("credhub", "api", "--server", config.ApiUrl, "--skip-tls-validation")
		RunCommand("credhub", "login", "--skip-tls-validation", "-u", config.ApiUsername, "-p", config.ApiPassword)

		CleanupCredhub(bbrTestPath)
	})

	AfterEach(func() {
		CleanupCredhub(bbrTestPath)
		CleanupArtifacts()
	})

	It("Successfully backs up and restores a Credhub release", func() {
		By("adding a test credential")
		RunCommand("credhub", "set", "--name", credentialName, "--type", "password", "-w", "originalsecret")

		By("running bbr backup")
		RunCommand("bbr", "director", "--private-key-path", config.Bosh.SshPrivateKeyPath,
			"--username", config.Bosh.SshUsername, "--host", config.Bosh.Host, "backup")

		By("asserting that the backup archive exists and contains a pg dump file")
		RunCommand("sh", "-c", fmt.Sprintf("tar -xvf ./%s*Z/bosh*credhub.tar", config.DirectorHost))
		Eventually(RunCommand("ls", "credhubdb_dump")).Should(gexec.Exit(0))

		By("editing the test credential")
		RunCommand("credhub", "set", "--name", credentialName, "--type", "password", "-w", "updatedsecret")

		editSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(editSession).Should(gexec.Exit(0))
		Eventually(editSession.Out).Should(gbytes.Say("value: updatedsecret"))

		By("running bbr restore")
		RunCommand("sh", "-c",
			fmt.Sprintf("bbr director --private-key-path %s --username %s --host %s restore --artifact-path ./%s*Z/",
				config.Bosh.SshPrivateKeyPath, config.Bosh.SshUsername, config.Bosh.Host,config.DirectorHost))

		By("checking if the test credentials was restored")
		getSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(getSession).Should(gexec.Exit(0))
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
