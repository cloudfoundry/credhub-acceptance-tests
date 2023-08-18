package integration_test

import (
	"io/ioutil"
	"os"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Import/Export test", func() {

	When("Importing an export file", func() {
		const (
			credentialRootPath               = "/bruce"
			selfSignedRootCACredPath         = credentialRootPath + "/ca"
			rootSignedLeafCredPath           = credentialRootPath + "/ca-leaf"
			rootSignedIntermediateCACredPath = credentialRootPath + "/int"
			intermediateSignedLeafCredPath   = credentialRootPath + "/int-leaf"
		)

		var (
			session    *Session
			exportFile *os.File
		)

		BeforeEach(func() {
			var err error
			exportFile, err = ioutil.TempFile("", "export-data")
			Expect(err).NotTo(HaveOccurred())

			session = RunCommand("generate",
				"--name", selfSignedRootCACredPath,
				"--type", "certificate",
				"--self-sign",
				"--is-ca",
				"--common-name", "bruce-ca",
			)
			Expect(session).To(Exit(0))

			session = RunCommand("generate",
				"--name", rootSignedIntermediateCACredPath,
				"--type", "certificate",
				"--ca", selfSignedRootCACredPath,
				"--is-ca",
				"--common-name", "bruce-int",
			)
			Expect(session).To(Exit(0))

			session = RunCommand("generate",
				"--name", rootSignedLeafCredPath,
				"--type", "certificate",
				"--ca", selfSignedRootCACredPath,
				"--common-name", "bruce-leaf1",
			)
			Expect(session).To(Exit(0))

			session = RunCommand("generate",
				"--name", intermediateSignedLeafCredPath,
				"--type", "certificate",
				"--ca", rootSignedIntermediateCACredPath,
				"--common-name", "bruce-leaf2",
			)
			Expect(session).To(Exit(0))
		})

		AfterEach(func() {
			session = RunCommand("delete",
				"--path", credentialRootPath,
			)
			Expect(session).To(Exit(0))
		})

		It("should restore the exported credentials", func() {
			selfSignedRootCA := getTrimmedCertificateForComparison(selfSignedRootCACredPath)
			rootSignedIntermediateCA := getTrimmedCertificateForComparison(rootSignedIntermediateCACredPath)
			rootSignedLeaf := getTrimmedCertificateForComparison(rootSignedLeafCredPath)
			intermediateSignedLeaf := getTrimmedCertificateForComparison(intermediateSignedLeafCredPath)

			session = RunCommand("export",
				"--path", credentialRootPath,
				"--file", exportFile.Name(),
			)
			Expect(session).To(Exit(0))

			session = RunCommand("delete",
				"--path", credentialRootPath,
			)
			Expect(session).To(Exit(0))

			session = RunCommand("import",
				"-f", exportFile.Name(),
			)
			Expect(session).To(Exit(0))

			Expect(getTrimmedCertificateForComparison(selfSignedRootCACredPath)).To(Equal(selfSignedRootCA))
			Expect(getTrimmedCertificateForComparison(rootSignedIntermediateCACredPath)).To(Equal(rootSignedIntermediateCA))
			Expect(getTrimmedCertificateForComparison(rootSignedLeafCredPath)).To(Equal(rootSignedLeaf))
			Expect(getTrimmedCertificateForComparison(intermediateSignedLeafCredPath)).To(Equal(intermediateSignedLeaf))
		})
	})

	When("Importing an export a self signed cert without a ca", func() {
		const (
			credentialRootPath = "/tobi"
			selfSignedCertPath = credentialRootPath + "/bruce-self-cert"
		)

		var (
			session    *Session
			exportFile *os.File
		)

		BeforeEach(func() {
			var err error
			exportFile, err = ioutil.TempFile("", "export-data")
			Expect(err).NotTo(HaveOccurred())

			session = RunCommand("generate",
				"--name", selfSignedCertPath,
				"--type", "certificate",
				"--self-sign",
				"--common-name", "bruce-ca",
			)
			Expect(session).To(Exit(0))
		})

		AfterEach(func() {
			session = RunCommand("delete",
				"--path", credentialRootPath,
			)
			Expect(session).To(Exit(0))
		})

		It("should restore the exported credentials", func() {
			selfSignedCert := getTrimmedCertificateForComparison(selfSignedCertPath)

			session = RunCommand("export",
				"--path", credentialRootPath,
				"--file", exportFile.Name(),
			)
			Expect(session).To(Exit(0))

			session = RunCommand("delete",
				"--path", credentialRootPath,
			)
			Expect(session).To(Exit(0))

			session = RunCommand("import",
				"-f", exportFile.Name(),
			)
			Expect(session).To(Exit(0))

			Expect(getTrimmedCertificateForComparison(selfSignedCertPath)).To(Equal(selfSignedCert))
		})
	})
})

func getTrimmedCertificateForComparison(name string) map[string]interface{} {
	session = RunCommand("get", "--name", name)
	Expect(session).To(Exit(0))

	certificate := make(map[string]interface{})
	err := yaml.Unmarshal(session.Out.Contents(), certificate)
	Expect(err).NotTo(HaveOccurred())
	delete(certificate, "version_created_at")
	delete(certificate, "id")

	return certificate
}
