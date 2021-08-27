package main

import (
	"log"
	"strings"
)

const (
	codeTtPath = "main_template.go.raw"
)

type solution struct {
	CodeTt string
}

func newSolution() solution {
	return solution{}
}

func main() {
	app := newSolution()

	err := checkErrors(
		app.getCodeTt, // get code tt
		app.updateCode,
		app.saveCode,
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("code generated successful!")
}

func (sol *solution) getCodeTt() error {
	var err error
	sol.CodeTt, err = ReadFileToString(codeTtPath)
	return err
}

func (sol *solution) updateCode() error {
	encryptionKey := getRandomString(32)
	sol.CodeTt = strings.ReplaceAll(sol.CodeTt, "%EncryptionKey%", encryptionKey)
	return nil
}

func (sol *solution) saveCode() error {
	// TODO
	return nil
}
