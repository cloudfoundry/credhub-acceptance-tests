package bbr_integration

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type JumpBoxSession struct {
	WorkspaceDir    string
	BinaryPath      string
	CertificatePath string
}

func NewJumpBoxSession(uniqueTestID string) *JumpBoxSession {
	session := JumpBoxSession{}
	session.WorkspaceDir = "/var/vcap/store/backup_workspace" + uniqueTestID
	var bbrBuildPath = MustHaveEnv("BBR_BUILD_PATH")

	By("setting up the jump box")
	RunOnJumpboxAsVcap("sudo mkdir", session.WorkspaceDir, "&& sudo chown vcap:vcap", session.WorkspaceDir, "&& sudo chmod 0777", session.WorkspaceDir)

	By("copying the BBR binary to the jumpbox")
	bbrFilename := getTar(bbrBuildPath)
	CopyOnJumpbox(
		filepath.Join(bbrBuildPath, bbrFilename),
		fmt.Sprintf("%s:%s", JumpboxInstance(), session.WorkspaceDir),
	)

	RunOnJumpboxAsVcap("tar -C", session.WorkspaceDir, "-xvf", filepath.Join(session.WorkspaceDir, bbrFilename))
	session.BinaryPath = filepath.Join(session.WorkspaceDir, "releases", "bbr")

	By("copying the Director certificate to the jumpbox")
	CopyOnJumpbox(
		MustHaveEnv("BOSH_CERT_PATH"),
		fmt.Sprintf("%s:%s/bosh.crt", JumpboxInstance(), session.WorkspaceDir),
	)
	session.CertificatePath = filepath.Join(session.WorkspaceDir, "bosh.crt")

	return &session
}

func (jumpBoxSession *JumpBoxSession) Cleanup() {
	By("remove workspace directory")
	RunOnJumpbox("sudo rm -rf", jumpBoxSession.WorkspaceDir)
}

func getTar(path string) string {
	glob := filepath.Join(path, "*.tar")
	matches, err := filepath.Glob(glob)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(1), "There should be only one tar file present in the BBR binary path")
	return filepath.Base(matches[0])
}
