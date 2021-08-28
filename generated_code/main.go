package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"time"
)

const (
	encryptFromFilePath = "encrypt_this.txt"
	encryptToFilePath   = "encrypted.dat"
	decryptToFilePath   = "decrypted.txt"
)

func main() {
	flagMode := flag.String("mode", "encrypt", "encrypt / decrypt")
	flagDuration := flag.Int("duration", 1, "duration amount")
	flagType := flag.String("type", "Y", "duration type: Y, M, D")

	// parse duration
	var durationAmount int = *flagDuration
	var lockDuration time.Duration
	switch *flagType {
	default:
		log.Fatalln("unknown duration type given (or not specified?)")
	case "Y":
		lockDuration = time.Hour * time.Duration(365*24*durationAmount)
	case "M":
		lockDuration = time.Hour * time.Duration(30*24*durationAmount)
	case "D":
		lockDuration = time.Hour * time.Duration(24*durationAmount)
	}

	switch *flagMode {
	default:
		log.Fatalln("unknown mode specified (or not specified?)")
	case "encrypt":
		encrypt(lockDuration)
		return
	case "decrypt":
		decrypt(lockDuration)
		return
	}
}

func encrypt(duration time.Duration) error {
	// read file
	fileBytes, err := readFile(encryptFromFilePath)
	if err != nil {
		return err
	}

	// create data container
	timeFrom := time.Now()
	timeTo := timeFrom.Add(duration)
	data := dataContainer{
		CreatedOn: timeFrom.UTC().String(),
		UnlockOn:  timeTo.UTC().String(),
		Data:      fileBytes,
	}

	// encode data container to json
	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return errors.New("failed to encode data container to json: " +
			err.Error())
	}

	// encrypt data container
	encryptedData, err := rsaEncrypt(jsonBytes, encryptionKey)
	if err != nil {
		return err
	}

	// write data container to file
	err = saveToFile(encryptToFilePath, encryptedData)
	if err != nil {
		return err
	}
	return nil
}

func decrypt(duration time.Duration) {
	// TODO
	//tHandler := newTimeHandler()
	//tHandler.parseTimeFromWorldAPI()
}
