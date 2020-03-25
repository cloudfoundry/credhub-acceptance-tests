package acceptance_test

import (
	"strings"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"strconv"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
)

var _ = Describe("Find", func() {
	currentTime := time.Now().UnixNano()
	randomString := time.Now().UnixNano()

	passwordName1 := testCredentialPath(randomString, fmt.Sprint("find-test/first-password-", currentTime))
	passwordName2 := testCredentialPath(randomString, fmt.Sprint("find-test/second-password-", currentTime))

	passwordPrefix := strings.SplitAfter(passwordName1, "find-test")[0]

	var expectedResult credentials.FindResults

	BeforeEach(func() {
		var err error

		generatePassword := generate.Password{Length: 10}

		expectedPassword1, err := credhubClient.GeneratePassword(passwordName1, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())

		expectedPassword2, err := credhubClient.GeneratePassword(passwordName2, generatePassword, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())

		expectedResult = credentials.FindResults{Credentials: []struct {
			Name             string `json:"name" yaml:"name"`
			VersionCreatedAt string `json:"version_created_at" yaml:"version_created_at"`
		}{
			{Name: passwordName2, VersionCreatedAt: expectedPassword2.VersionCreatedAt},
			{Name: passwordName1, VersionCreatedAt: expectedPassword1.VersionCreatedAt},
		}}
	})

	AfterEach(func() {
		err := credhubClient.Delete(passwordName1)
		Expect(err).ToNot(HaveOccurred())
		err = credhubClient.Delete(passwordName2)
		Expect(err).ToNot(HaveOccurred())
	})

	Specify("finding the credentials by path", func() {
		results, err := credhubClient.FindByPath(passwordPrefix)
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(Equal(expectedResult))
	})

	Specify("finding the credentials by name-like", func() {
		results, err := credhubClient.FindByPartialName(strconv.FormatInt(currentTime, 10))
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(Equal(expectedResult))
	})
})
