package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
)

var _ = Describe("Import test", func() {
	It("should import credentials from a file", func() {
		RunCommand("generate", "-n", "ca-certificate", "-t", "certificate", "-c", "credhub-ca", "-o", "pivotal", "-u", "credhub", "-i", "nyc", "-s", "NY", "-y", "US", "--is-ca", "--self-sign")

		session := RunCommand("import", "-f", "../test_helpers/bulk_import.yml")
		Eventually(session).Should(Exit(0))

		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - agent`))
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - director`))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-ca`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-cert`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
	})

	It("should save the credentials on CredHub", func() {
		RunCommand("generate", "-n", "ca-certificate", "-t", "certificate", "-c", "credhub-ca", "-o", "pivotal", "-u", "credhub", "-i", "nyc", "-s", "NY", "-y", "US", "--is-ca", "--self-sign")
		RunCommand("import", "-f", "../test_helpers/bulk_import.yml")

		session := RunCommand("get", "-n", "/director/deployment/blobstore - agent")
		Eventually(session).Should(Exit(0))
		stdOut := string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - agent`))
		Expect(stdOut).To(ContainSubstring(`type: password`))
		Expect(stdOut).To(ContainSubstring(`value: gx4ll8193j5rw0wljgqo`))

		session = RunCommand("get", "-n", "/director/deployment/blobstore - director")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/blobstore - director`))
		Expect(stdOut).To(ContainSubstring(`type: value`))
		Expect(stdOut).To(ContainSubstring(`value: y14ck84ef51dnchgk4kp`))

		session = RunCommand("get", "-n", "/director/deployment/bosh-ca")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-ca`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))

		session = RunCommand("get", "-n", "/director/deployment/bosh-cert")
		Eventually(session).Should(Exit(0))
		stdOut = string(session.Out.Contents())
		Expect(stdOut).To(ContainSubstring(`name: /director/deployment/bosh-cert`))
		Expect(stdOut).To(ContainSubstring(`type: certificate`))
		Expect(stdOut).To(ContainSubstring(`value:`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
		Expect(stdOut).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
	})

})
