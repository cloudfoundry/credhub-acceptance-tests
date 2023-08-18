package integration

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	config Config
	err    error
)

var _ = Describe("vcap interpolation of secrets", func() {
	credentialName := GenerateUniqueCredentialName()
	credentialValue := `{"username":"bob", "password":"bob has a password"}`

	BeforeEach(func() {
		config, err = LoadConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should interpolate a known secret into the VCAP_SERVICES json", func() {
		By("setting a new json secret", func() {
			session := RunCommand("set", "-n", credentialName, "-t", "json", "-v", credentialValue)
			Eventually(session).Should(Exit(0))
		})

		By("posting the VCAP_SERVICES JSON", func() {
			session := RunCommand("--token")
			Eventually(session).Should(Exit(0))
			token := strings.TrimSpace(string(session.Out.Contents()))

			postData := `{` +
				`"p-config-server": [` +
				`   {` +
				`     "credentials": {` +
				`       "credhub-ref": "((/` + credentialName + `))"` +
				`     },` +
				`     "label": "p-config-server"` +
				`   }` +
				` ]` +
				`}`

			result, statusCode, err := postJSON(config.ApiUrl+"/api/v1/interpolate", postData, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))

			Expect(result).To(Equal(`{"p-config-server":[{"credentials":{"password":"bob has a password","username":"bob"},"label":"p-config-server"}]}`))
		})
	})
})

func postJSON(url string, postData string, token string) (string, int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, strings.NewReader(postData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(body), resp.StatusCode, nil
}
