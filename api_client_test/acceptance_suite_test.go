package acceptance_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
)

var (
	credhubClient *credhub.CredHub
)

var _ = BeforeEach(func() {
	var err error

	config, err := LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	credhub_ca, err := ioutil.ReadFile(path.Join(config.CredentialRoot, "server_ca_cert.pem"))
	Expect(err).NotTo(HaveOccurred())

	uaa_ca, err := ioutil.ReadFile(path.Join(config.UAACa))
	Expect(err).NotTo(HaveOccurred())

	credhubClient, err = credhub.New(config.ApiUrl,
		credhub.CaCerts(string(credhub_ca), string(uaa_ca)),
		credhub.Auth(
			auth.UaaClientCredentials(config.ClientName, config.ClientSecret),
		),
	)

	Expect(err).ToNot(HaveOccurred())
})

func TestCredhub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Client Suite")
}

func testCredentialPath(randomizer int64, credentialName string) string {
	return fmt.Sprintf("/acceptance/%v/%v", randomizer, credentialName)
}
