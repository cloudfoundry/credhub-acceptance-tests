package acceptance_test

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

var _ = Describe("SSH Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath(time.Now().UnixNano(), "some-ssh")
		generateParameters := generate.SSH{KeyLength: 2048}

		By("generate ssh keys with path " + name)
		generatedSSH, err := credhubClient.GenerateSSH(name, generateParameters, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		block, _ := pem.Decode([]byte(generatedSSH.Value.PrivateKey))
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(privateKey.N.BitLen()).To(Equal(generateParameters.KeyLength))

		By("generate the ssh keys again without overwrite returns same ssh")
		ssh, err := credhubClient.GenerateSSH(name, generateParameters, credhub.NoOverwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh).To(Equal(generatedSSH))

		By("overwriting with generate")
		ssh, err = credhubClient.GenerateSSH(name, generateParameters, credhub.Overwrite)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh).ToNot(Equal(generatedSSH))

		By("setting the ssh keys again overwrites previous ssh")
		newSSH := values.SSH{PrivateKey: "private key", PublicKey: "public key"}
		ssh, err = credhubClient.SetSSH(name, newSSH)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh.Value.SSH).To(Equal(newSSH))

		By("getting the ssh credential")
		ssh, err = credhubClient.GetLatestSSH(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh.Value.SSH).To(Equal(newSSH))

		By("deleting the rsa credential")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetLatestRSA(name)
		Expect(err).To(HaveOccurred())
	})
})
