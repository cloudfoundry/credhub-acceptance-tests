package integration_test

import (
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	credentialNamesSet []string
	credentialNamesGet []string
	session            *Session
)

var _ = Describe("Import test", func() {

	BeforeEach(func() {
		session = RunCommand("set", "-n", "ca-certificate", "-t", "certificate", "-c", VALID_CERTIFICATE_CA)
		Eventually(session).Should(Exit(0))
	})

	It("should import credentials from a file", func() {
		session = RunCommand("import", "-f", "../test_helpers/bulk_import_set.yml")
		Eventually(session).Should(Exit(0))

		credentialNamesSet = []string{
			"/director/deployment/blobstore-agent1",
			"/director/deployment/blobstore-director1",
			"/director/deployment/bosh-ca1",
			"/director/deployment/bosh-cert1",
			"/director/deployment/ssh-cred1",
			"/director/deployment/rsa-cred1",
			"/director/deployment/user1",
			"/director/deployment/json1",
		}

		for _, credentialName := range credentialNamesSet {
			session = RunCommand("delete", "-n", credentialName)
			Eventually(session).Should(Exit(0))
		}
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

	Describe("when there is a cert chain and ca_name is given for a cert and the ca is later in the file", func() {
		It("should import the ca before the cert it signed", func() {
			beforeCertChainGet()

			session = RunCommand("curl", "-p", "api/v1/certificates?name=/leaf_cert")
			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`"signed_by": "/intermediate_ca"`))

			session = RunCommand("curl", "-p", "api/v1/certificates?name=/intermediate_ca")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`"signed_by": "/root_ca"`))
			Expect(stdOut).To(ContainSubstring("leaf_cert"))

			session = RunCommand("curl", "-p", "api/v1/certificates?name=/root_ca")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring("intermediate_ca"))

			afterGet()
		})
	})

	Describe("when the credentials have metadata", func() {
		BeforeEach(func() {
			supported, err := serverSupportsMetadata()
			Expect(err).NotTo(HaveOccurred())
			if !supported {
				Skip("Server does not support metadata")
			}
		})

		It("should import credentials from a file", func() {
			session = RunCommand("import", "-f", "../test_helpers/bulk_import_set_with_metadata.yml")
			Eventually(session).Should(Exit(0))

			session = RunCommand("get", "-n", "/director/deployment/blobstore-agent1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut := string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore-agent1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: password`))
			Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some: metadata`))

			session = RunCommand("get", "-n", "/director/deployment/blobstore-director1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore-director1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: value`))
			Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some:
  - different
  - metadata`))

			session = RunCommand("get", "-n", "/director/deployment/bosh-ca1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-ca1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: certificate`))
			Expect(stdOut).To(ContainSubstring(`value:`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some:
    object:
      with: data`))

			session = RunCommand("get", "-n", "/director/deployment/bosh-cert1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-cert1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: certificate`))
			Expect(stdOut).To(ContainSubstring(`value:`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some:
    object:
    - with
    - an
    - array`))

			session = RunCommand("get", "-n", "/director/deployment/ssh-cred1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/ssh-cred1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: ssh`))
			Expect(stdOut).To(ContainSubstring(`value:`))
			Expect(stdOut).To(ContainSubstring(`ssh-rsa`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  top:
    one: foo
    two: bar`))

			session = RunCommand("get", "-n", "/director/deployment/rsa-cred1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/rsa-cred1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: rsa`))
			Expect(stdOut).To(ContainSubstring(`value:`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN PUBLIC KEY-----`))
			Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some: thing`))

			session = RunCommand("get", "-n", "/director/deployment/user1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/user1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: user`))
			Expect(stdOut).To(ContainSubstring(`value:`))
			Expect(stdOut).To(ContainSubstring(`password: lGcaYF31nJNCii53OkNhtjo9tXJ3kf`))
			Expect(stdOut).To(ContainSubstring(`username: dan-user1`))
			Expect(stdOut).To(ContainSubstring(`password_hash:`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some: thing`))

			session = RunCommand("get", "-n", "/director/deployment/json1-with-metadata")
			Eventually(session).Should(Exit(0))
			stdOut = string(session.Out.Contents())
			Expect(stdOut).To(ContainSubstring(`name: /director/deployment/json1-with-metadata`))
			Expect(stdOut).To(ContainSubstring(`type: json`))
			Expect(stdOut).To(ContainSubstring(`value:`))
			Expect(stdOut).To(ContainSubstring(`trump:`))
			Expect(stdOut).To(ContainSubstring(`tweet:`))
			Expect(stdOut).To(ContainSubstring(`- covfefe`))
			Expect(stdOut).To(ContainSubstring(`- covfefe`))
			Expect(stdOut).To(ContainSubstring(`
metadata:
  some: thing`))

			credentialNamesSet = []string{
				"/director/deployment/blobstore-agent1-with-metadata",
				"/director/deployment/blobstore-director1-with-metadata",
				"/director/deployment/bosh-ca1-with-metadata",
				"/director/deployment/bosh-cert1-with-metadata",
				"/director/deployment/ssh-cred1-with-metadata",
				"/director/deployment/rsa-cred1-with-metadata",
				"/director/deployment/user1-with-metadata",
				"/director/deployment/json1-with-metadata",
			}

			for _, credentialName := range credentialNamesSet {
				session = RunCommand("delete", "-n", credentialName)
				Eventually(session).Should(Exit(0))

			}
		})
	})
})

func beforeGet() {
	session = RunCommand("import", "-f", "../test_helpers/bulk_import_get.yml")
	Eventually(session).Should(Exit(0))
	credentialNamesGet = []string{
		"/director/deployment/blobstore-agent2",
		"/director/deployment/blobstore-director2",
		"/director/deployment/bosh-ca2",
		"/director/deployment/bosh-cert2",
		"/director/deployment/ssh-cred2",
		"/director/deployment/rsa-cred2",
		"/director/deployment/user2",
		"/director/deployment/json2",
	}
}

func beforeCertChainGet() {
	session = RunCommand("import", "-f", "../test_helpers/bulk_import_with_ca_name.yml")
	Eventually(session).Should(Exit(0))
	credentialNamesGet = []string{
		"root_ca",
		"intermediate_ca",
		"leaf_cert",
	}
}

func afterGet() {
	for _, credentialName := range credentialNamesGet {
		session = RunCommand("delete", "-n", credentialName)
		Eventually(session).Should(Exit(0))
	}
}
