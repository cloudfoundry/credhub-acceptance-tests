package utilities

import (
	"fmt"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SessionInput struct {
	fullCommand string
	commandName string
}

type CliCommand struct {
	sessionInput  SessionInput
	sessionOutput string
}

func GenerateAutoDoc(session *gexec.Session) error {
	sessionOutput := string(session.Out.Contents())
	sessionInput := parseSessionInput(session.Command.Args)
	cliCommand := CliCommand{
		sessionInput:  sessionInput,
		sessionOutput: sessionOutput,
	}
	err := writeFiles(cliCommand)
	return err
}

func parseSessionInput(args []string) SessionInput {
	args[0] = "credhub"
	commandName := args[1]
	fullCommand := strings.Join(args, " ")

	return SessionInput{fullCommand, commandName}
}

func writeFiles(command CliCommand) error {
	folderName := command.sessionInput.commandName
	folderPath := filepath.Join("/tmp/credhub_cli_docs", folderName)
	err := createDirectory(folderPath)
	if err != nil {
		return err
	}
	inputFilePath := filepath.Join(folderPath, "input.adoc")
	err = generateInputFile(command.sessionInput, inputFilePath)
	if err != nil {
		return err
	}
	outputFilePath := filepath.Join(folderPath, "output.adoc")
	err = generateOutputFile(command.sessionOutput, outputFilePath)
	if err != nil {
		return err
	}

	return nil

}

func generateInputFile(input SessionInput, path string) error {
	formattedInput := fmt.Sprintf("```\n" + input.fullCommand + "\n" + "```")
	err := ioutil.WriteFile(path, []byte(formattedInput), os.ModePerm)
	return err

}

func generateOutputFile(output string, path string) error {
	formattedOutput := fmt.Sprintf("```\n" + output + "\n" + "```")
	err := ioutil.WriteFile(path, []byte(formattedOutput), os.ModePerm)
	return err
}

func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}
