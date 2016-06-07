package commands_test

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"os/exec"
	"testing"
)

var commandPath string

var _ = Describe("Integration test", func() {
	It("smoke tests ok", func() {
		session := runCommand("api", "pivotal-credential-manager.cfapps.io")
		Eventually(session).Should(Exit(0))
		uniqueId := strconv.FormatInt(time.Now().UnixNano(), 10)
		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(1))
		session = runCommand("set", "-n", uniqueId, "-v", "bar")
		Eventually(session).Should(Exit(0))
		session = runCommand("get", "-n", uniqueId)
		Eventually(session).Should(Exit(0))
		fmt.Println(string(session.Out.Contents()))
		session = runCommand("delete", "-n", uniqueId)
		Eventually(session).Should(Exit(0))

	})
})

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/pivotal-cf/cm-cli")
	Expect(err).NotTo(HaveOccurred())
	return []byte(path)
}, func(data []byte) {
	commandPath = string(data)
})


func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}
