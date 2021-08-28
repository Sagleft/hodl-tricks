package main

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	codeTtPath          = "consts.go.raw"
	destinationCodePath = "generated_code/consts.go"
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

	log.Println("new encryption key generated!")
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
	return SaveStringToFile(destinationCodePath, sol.CodeTt)
}

// ReadFileToString read file to string
func ReadFileToString(filepath string) (string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", errors.New("failed to read file: " + err.Error())
	}
	return string(file), nil
}

type errorFunc func() error

func checkErrors(errChecks ...errorFunc) error {
	for _, errFunc := range errChecks {
		err := errFunc()
		if err != nil {
			return err
		}
	}
	return nil
}

func getRandomString(length int) string {
	rand.Seed(int64(time.Now().Nanosecond()))
	var symbolRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = symbolRunes[rand.Intn(len(symbolRunes))]
	}
	return string(b)
}

// SaveStringToFile save arbitrary string to file
func SaveStringToFile(filepath string, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return errors.New("failed to create file: " + err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return errors.New("failed to write string to file: " + err.Error())
	}
	return nil
}
