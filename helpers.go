package main

import "io/ioutil"

// ReadFileToString read file to string
func ReadFileToString(filepath string) (string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(file), err
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
