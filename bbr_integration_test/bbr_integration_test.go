package bbr_integration

import (
	"fmt"

	"os/exec"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Backup and Restore", func() {
	var credentialName string
	var bbrTestPath = "bbr_test"

	BeforeEach(func() {
		credentialName = fmt.Sprintf("%s/%s", bbrTestPath, GenerateUniqueCredentialName())

		By("authenticating against credhub")
		Eventually(Run(fmt.Sprintf(
			"%s api --server %s --skip-tls-validation; %s login --skip-tls-validation -u %s -p %s",
			credhubCliBinaryPath,
			credhubUrl,
			credhubCliBinaryPath,
			credhubApiUsername,
			credhubApiPassword,
		))).Should(gexec.Exit(0))

		CleanupCredhub(bbrTestPath)
	})

	AfterEach(func() {
		CleanupCredhub(bbrTestPath)
	})

	It("Successfully backs up and restores a Credhub release", func() {
		By("adding a test credential")
		Eventually(Run(fmt.Sprintf(
			"%s set --name %s --value originalsecret",
			credhubCliBinaryPath,
			credentialName,
		))).Should(gexec.Exit(0))

		By("running bbr backup")
		Eventually(Run(fmt.Sprintf(
			"cd %s; %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
			tmpDir,
			bbrBinaryPath,
			MustHaveEnv("BOSH_URL"),
			boshCertPath,
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			deploymentName,
		))).Should(gexec.Exit(0))

		By("asserting that the backup archive exists and contains a pg dump file")
		Eventually(Run(fmt.Sprintf(
			"cd %s/%s; tar zxvf %s; [ -f %s ]",
			tmpDir,
			deploymentName,
			"credhub-0.tgz",
			"./credhub/credhubdb_dump",
		))).Should(gexec.Exit(0))

		By("editing the test credential")
		Eventually(Run(fmt.Sprintf(
			"%s set --name %s --value updatedsecret",
			credhubCliBinaryPath,
			credentialName,
		))).Should(gexec.Exit(0))

		By("running bbr restore")
		Eventually(Run(fmt.Sprintf(
			"cd %s; %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s restore",
			tmpDir,
			bbrBinaryPath,
			MustHaveEnv("BOSH_URL"),
			boshCertPath,
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			deploymentName,
		))).Should(gexec.Exit(0))

		By("checking if the test credentials was restored")
		getSession := Run(fmt.Sprintf(
			"%s get --name %s",
			credhubCliBinaryPath,
			credentialName,
		))
		Eventually(getSession).Should(gexec.Exit(0))
		Eventually(getSession.Out).Should(gbytes.Say("value: originalsecret"))
	})
})

func CleanupCredhub(path string) {
	By("Cleaning up credhub bbr test passwords")
	Eventually(Run(fmt.Sprintf(
		"%s find -p /%s | tail -n +2 | cut -d\" \" -f1 | xargs -IN %s delete --name N",
		credhubCliBinaryPath,
		path,
		credhubCliBinaryPath,
	))).Should(gexec.Exit(0))
}

func Run(command string) *gexec.Session {
	session, err := gexec.Start(exec.Command("sh", "-c", command), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
