package utilities

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var _ = Describe("Auto Doc Generator Test", func() {
	Context("Generate a doc", func() {
		Context("parses session input", func() {
			It("replaces credhub path with credhub and joins array to string", func() {
				sessionInput := parseSessionInput([]string{"path/to/credhub-cli", "some-command", "some-arg"})
				Expect(sessionInput.fullCommand).To(Equal("credhub some-command some-arg"))
			})
			It("parses function name", func() {
				sessionInput := parseSessionInput([]string{"path/to/credhub-cli", "some-command", "some-arg"})
				Expect(sessionInput.commandName).To(Equal("some-command"))
			})

		})
		Context("generates formatted files", func() {
			var (
				path         string
				sessionInput SessionInput
				session      gexec.Session
			)

			BeforeEach(func() {
				session = gexec.Session{
					Command: &exec.Cmd{
						Path: "/var/folders/c4/nd0g0tkn10zcf19r5tsjzlwh0000gn/T/gexec_artifacts663775029/g605623568/credhub-cli",
						Args: []string{
							"/var/folders/c4/nd0g0tkn10zcf19r5tsjzlwh0000gn/T/gexec_artifacts663775029/g605623568/credhub-cli",
							"test-function",
							"-a",
							"actor-1551712691979175000",
							"-p",
							"/path-1551712691979175000",
							"-o",
							"read,write",
						},
					},
					Out: gbytes.BufferWithBytes([]byte("test")),
				}
				sessionInput = parseSessionInput(session.Command.Args)
				path = filepath.Join("/tmp/credhub_cli_docs", sessionInput.commandName)
				_ = os.RemoveAll(path)
			})

			AfterEach(func() {
				_ = os.RemoveAll(path)
			})

			It("properly generates input and output files", func() {
				err := GenerateAutoDoc(&session)
				Expect(err).NotTo(HaveOccurred())

				Expect(path).To(BeADirectory())

				inputFilePath := filepath.Join(path, "input.adoc")
				Expect(inputFilePath).To(BeAnExistingFile())
				actualInputFileData, err := ioutil.ReadFile(inputFilePath)
				expectedInputFileData := []byte(fmt.Sprintf("```\n" + sessionInput.fullCommand + "\n" + "```"))
				Expect(actualInputFileData).To(Equal(expectedInputFileData))
				Expect(err).NotTo(HaveOccurred())

				outputFilePath := filepath.Join(path, "output.adoc")
				Expect(outputFilePath).To(BeAnExistingFile())
				actualOutputFileData, err := ioutil.ReadFile(outputFilePath)
				expectedOutputFileData := []byte(fmt.Sprintf("```\n" + string(session.Out.Contents()) + "\n" + "```"))
				Expect(actualOutputFileData).To(Equal(expectedOutputFileData))
				Expect(err).NotTo(HaveOccurred())

			})
		})
	})
})
