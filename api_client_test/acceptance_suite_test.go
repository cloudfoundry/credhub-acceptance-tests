package acceptance_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

var (
	currentTestNumber = time.Now().UnixNano()
	credhubClient     *credhub.CredHub
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

func testCredentialPath(credentialName string) string {
	return fmt.Sprintf("/acceptance/%v/%v", currentTestNumber, credentialName)
}

func testCredentialPrefix() string {
	return fmt.Sprintf("/acceptance/%v/", currentTestNumber)
}
