package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"encoding/json"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

var _ = Describe("JSON Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-json")

		cred := make(map[string]interface{})
		cred["key"] = "value"

		cred2 := make(map[string]interface{})
		cred2["another_key"] = "another_value"

		var unmarshalled values.JSON

		By("setting the json for the first time returns same json")
		jsonValue, err := credhubClient.SetJSON(name, cred, credhub.Overwrite)

		keyValueJson := `{"key":"value"}`

		json.Unmarshal([]byte(keyValueJson), &unmarshalled)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(unmarshalled))

		By("setting the json again without overwrite returns same json")
		jsonValue, err = credhubClient.SetJSON(name, cred2, credhub.NoOverwrite)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(unmarshalled))

		By("overwriting the json with set")
		jsonValue, err = credhubClient.SetJSON(name, cred2, credhub.Overwrite)

		otherKeyValueJson := `{"another_key":"another_value"}`

		unmarshalled = nil

		json.Unmarshal([]byte(otherKeyValueJson), &unmarshalled)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(unmarshalled))

		By("getting the json")
		jsonValue, err = credhubClient.GetLatestJSON(name)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(unmarshalled))

		By("deleting the json")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestJSON(name)
		Expect(err).To(HaveOccurred())
	})
})
