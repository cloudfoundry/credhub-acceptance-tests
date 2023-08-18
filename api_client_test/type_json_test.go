package acceptance_test

import (
	"encoding/json"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-json")

		cred := make(map[string]interface{})
		cred["key"] = "value"

		cred2 := make(map[string]interface{})
		cred2["another_key"] = "another_value"

		var firstUnmarshalled values.JSON
		var secondUnmarshalled values.JSON

		By("setting the json for the first time returns same json")
		jsonValue, err := credhubClient.SetJSON(name, cred)

		json.Unmarshal([]byte(`{"key":"value"}`), &firstUnmarshalled)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(firstUnmarshalled))

		By("setting the json again overwrites previous json")
		jsonValue, err = credhubClient.SetJSON(name, cred2)

		json.Unmarshal([]byte(`{"another_key":"another_value"}`), &secondUnmarshalled)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(secondUnmarshalled))

		By("getting the json")
		jsonValue, err = credhubClient.GetLatestJSON(name)

		Expect(err).ToNot(HaveOccurred())
		Expect(jsonValue.Value).To(Equal(secondUnmarshalled))

		By("deleting the json")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestJSON(name)
		Expect(err).To(HaveOccurred())
	})
})
