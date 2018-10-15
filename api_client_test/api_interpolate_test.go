package acceptance_test

import (
	"encoding/json"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InterpolateString", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("thing-to-interpolate")
		cred := make(map[string]interface{})
		cred["firstKey"] = "this is a value for the first key"

		var unmarshalled values.JSON

		_, err := credhubClient.SetJSON(name, cred)
		Expect(err).ToNot(HaveOccurred())
		keyValueJson := `{"firstKey":"this is a value for the first key"}`
		json.Unmarshal([]byte(keyValueJson), &unmarshalled)

		By("calling interpolateString with a valid credhub-ref")
		vcapServicesString := `{"someOuterKey":[{"credentials": {"credhub-ref":"((` + name + `))"},"anotherKey":"some other value that does not change"}]}`
		interpolatedString, err := credhubClient.InterpolateString(vcapServicesString)
		Expect(err).ToNot(HaveOccurred())
		Expect(interpolatedString).To(Equal(`{"someOuterKey":[{"anotherKey":"some other value that does not change","credentials":{"firstKey":"this is a value for the first key"}}]}`))
		Expect(interpolatedString).NotTo(ContainSubstring("credhub-ref"))

		By("calling interpolateString without credential angd credhub-ref returns input unchanged")
		invalidString := `{"someOuterKey":[{"things":{"keyTwo":"this is not from credhub"},"anotherKey":"some other value that does not change"}]}`
		interpolatedString, err = credhubClient.InterpolateString(invalidString)
		Expect(err).ToNot(HaveOccurred())
		Expect(interpolatedString).To(Equal(invalidString))

		By("clean up from test")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestJSON(name)
		Expect(err).To(HaveOccurred())
	})
})
