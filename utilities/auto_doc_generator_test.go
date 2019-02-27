package utilities

import (
	"fmt"
	. "github.com/onsi/ginkgo"
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
			session := gexec.Session{
				Command: &exec.Cmd{
					Path: "/var/folders/c4/nd0g0tkn10zcf19r5tsjzlwh0000gn/T/gexec_artifacts663775029/g605623568/credhub-cli",
					Args: []string {
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
			sessionInput := parseSessionInput(session.Command.Args)
			path := filepath.Join("/tmp/credhub_cli_docs", sessionInput.commandName)

			BeforeEach(func() {
				_ = os.RemoveAll(path)
			})

			AfterSuite(func() {
				_ = os.RemoveAll(path)
			})

			It("creates output folder if not exists", func() {
				err := GenerateAutoDoc(&session)

				Expect(err).NotTo(HaveOccurred())
				Expect(path).To(BeADirectory())
			})
			It("creates input file with session input data", func() {
				err := GenerateAutoDoc(&session)

				inputFilePath := filepath.Join(path, "input.adoc")
				Expect(inputFilePath).To(BeAnExistingFile())

				actualFileData, err := ioutil.ReadFile(inputFilePath)
				expectedFileData := []byte(fmt.Sprintf("```\n" + sessionInput.fullCommand + "\n" + "```"))

				Expect(actualFileData).To(Equal(expectedFileData))
				Expect(err).NotTo(HaveOccurred())
			})
			It("creates output file with session output data", func() {
				err := GenerateAutoDoc(&session)

				outputFilePath := filepath.Join(path, "output.adoc")
				Expect(outputFilePath).To(BeAnExistingFile())

				actualFileData, err := ioutil.ReadFile(outputFilePath)
				expectedFileData := []byte(fmt.Sprintf("```\n" + string(session.Out.Contents()) + "\n" + "```"))

				Expect(actualFileData).To(Equal(expectedFileData))
				Expect(err).NotTo(HaveOccurred())

			})
		})
	})
})
