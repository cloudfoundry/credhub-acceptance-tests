package bbr_integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Backup and Restore", func() {
	It("Creates a backup archive", func() {
		By("running bbr backup")
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf(
			"cd %s; %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
			jumpBoxSession.WorkspaceDir,
			jumpBoxSession.BinaryPath,
			MustHaveEnv("BOSH_URL"),
			jumpBoxSession.CertificatePath,
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			DeploymentToBackup(),
		))).Should(gexec.Exit(0))

		By("asserting that the backup archive exists and contains a pg dump file")
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf(
			"cd %s/%s; tar zxvf %s; [ -f %s ]",
			jumpBoxSession.WorkspaceDir,
			DeploymentToBackup(),
			"credhub-0.tgz",
			"./credhub/credhubdb_dump",
		))).Should(gexec.Exit(0))
	})
})
