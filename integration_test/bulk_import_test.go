package integration

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	credentialNamesSet []string
	credentialNamesGet []string
	session            *Session
)

var _ = Describe("Import test", func() {

	It("should import credentials from a file", func() {
		beforeSet()

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore-agent1`))
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore-director1`))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-ca1`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-cert1`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/ssh-cred1`))
		Expect(stdOut).To(ContainSubstring(`type: ssh`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`ssh-rsa`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/rsa-cred1`))
		Expect(stdOut).To(ContainSubstring(`type: rsa`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN PUBLIC KEY-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/user1`))
		Expect(stdOut).To(ContainSubstring(`type: user`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`password: lGcaYF31nJNCii53OkNhtjo9tXJ3kf`))
		Expect(stdOut).To(ContainSubstring(`username: dan-user1`))
		Expect(stdOut).To(ContainSubstring(`password_hash:`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/json1`))
		Expect(stdOut).To(ContainSubstring(`type: json`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`trump:`))
		Expect(stdOut).To(ContainSubstring(`tweet:`))
		Expect(stdOut).To(ContainSubstring(`- covfefe`))
		Expect(stdOut).To(ContainSubstring(`- covfefe`))
		Expect(stdOut).To(ContainSubstring(`- covfefe`))

		afterSet()
	})

	It("should save the credentials on CredHub", func() {
		beforeGet()

		session = RunCommand("get", "-n", "/director/deployment/blobstore-agent2")
		Eventually(session).Should(Exit(0))
		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore-agent2`))
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))

		session = RunCommand("get", "-n", "/director/deployment/blobstore-director2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore-director2`))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))

		session = RunCommand("get", "-n", "/director/deployment/bosh-ca2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-ca2`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		session = RunCommand("get", "-n", "/director/deployment/bosh-cert2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-cert2`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		session = RunCommand("get", "-n", "/director/deployment/ssh-cred2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/ssh-cred2`))
		Expect(stdOut).To(ContainSubstring(`type: ssh`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`ssh-rsa`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		session = RunCommand("get", "-n", "/director/deployment/rsa-cred2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/rsa-cred2`))
		Expect(stdOut).To(ContainSubstring(`type: rsa`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN PUBLIC KEY-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		session = RunCommand("get", "-n", "/director/deployment/user2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/user2`))
		Expect(stdOut).To(ContainSubstring(`type: user`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`password: lGcaYF31nJNCii53OkNhtjo9tXJ3kf`))
		Expect(stdOut).To(ContainSubstring(`username: dan-user2`))
		Expect(stdOut).To(ContainSubstring(`password_hash:`))

		session = RunCommand("get", "-n", "/director/deployment/json2")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/json2`))
		Expect(stdOut).To(ContainSubstring(`type: json`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`trump:`))
		Expect(stdOut).To(ContainSubstring(`tweet:`))
		Expect(stdOut).To(ContainSubstring(`- covfefe`))
		Expect(stdOut).To(ContainSubstring(`- covfefe`))
		Expect(stdOut).To(ContainSubstring(`- covfefe`))

		afterGet()
	})

})

func beforeSet() {
	session = RunCommand("set", "-n", "ca-certificate", "-t", "certificate", "-c", VALID_CERTIFICATE_CA)
	Eventually(session).Should(Exit(0))
	credentialNamesSet = []string{"ca-certificate"}
	session = RunCommand("import", "-f", "../test_helpers/bulk_import_set.yml")
	Eventually(session).Should(Exit(0))
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/blobstore-agent1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/blobstore-director1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/bosh-ca1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/bosh-cert1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/ssh-cred1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/rsa-cred1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/user1")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/json1")
}

func afterSet() {
	for _, credentialName := range credentialNamesSet {
		session = RunCommand("delete", "-n", credentialName)
		Eventually(session).Should(Exit(0))
	}
}

func beforeGet() {
	session = RunCommand("set", "-n", "ca-certificate", "-t", "certificate", "-c", VALID_CERTIFICATE_CA)
	Eventually(session).Should(Exit(0))
	credentialNamesGet = []string{"ca-certificate"}
	session = RunCommand("import", "-f", "../test_helpers/bulk_import_get.yml")
	Eventually(session).Should(Exit(0))
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/blobstore-agent2")
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/blobstore-director2")
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/bosh-ca2")
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/bosh-cert2")
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/ssh-cred2")
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/rsa-cred2")
	credentialNamesGet = append(credentialNamesGet, "/director/deployment/user2")
	credentialNamesSet = append(credentialNamesSet, "/director/deployment/json2")
}

func afterGet() {
	for _, credentialName := range credentialNamesGet {
		session = RunCommand("delete", "-n", credentialName)
		Eventually(session).Should(Exit(0))
	}
}
