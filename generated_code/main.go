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
	flag.Parse()

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
		result, err := encrypt(lockDuration)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("file encrypted to " + result.TimeTo.String())
		return
	case "decrypt":
		err := decrypt(lockDuration)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
}

func encrypt(duration time.Duration) (*encryptResult, error) {
	// read file
	fileBytes, err := readFile(encryptFromFilePath)
	if err != nil {
		return nil, err
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
		return nil, errors.New("failed to encode data container to json: " +
			err.Error())
	}

	// encrypt data container
	encryptedData, err := rsaEncrypt(jsonBytes, encryptionKey)
	if err != nil {
		return nil, err
	}

	// write data container to file
	err = saveToFile(encryptToFilePath, encryptedData)
	if err != nil {
		return nil, err
	}
	return &encryptResult{
		TimeTo: timeTo,
	}, nil
}

func decrypt(duration time.Duration) error {
	// TODO
	//tHandler := newTimeHandler()
	//tHandler.parseTimeFromWorldAPI()
	return nil
}
