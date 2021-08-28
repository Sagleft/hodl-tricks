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
	flagType := flag.String("type", "year", "duration type: year, month, day, hour, minute, second")
	flag.Parse()

	switch *flagMode {
	default:
		log.Fatalln("unknown mode specified (or not specified?)")
	case "encrypt":
		// parse duration
		var durationAmount int = *flagDuration
		var lockDuration time.Duration
		switch *flagType {
		default:
			log.Fatalln("unknown duration type given (or not specified?)")
		case "year":
			lockDuration = time.Hour * time.Duration(365*24*durationAmount)
		case "month":
			lockDuration = time.Hour * time.Duration(30*24*durationAmount)
		case "day":
			lockDuration = time.Hour * time.Duration(24*durationAmount)
		case "hour":
			lockDuration = time.Hour * time.Duration(durationAmount)
		case "minute":
			lockDuration = time.Minute * time.Duration(durationAmount)
		case "second":
			lockDuration = time.Second * time.Duration(durationAmount)
		}

		result, err := encrypt(lockDuration)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("[DONE] File encrypted to " + result.TimeTo.String())
		log.Println("(!!) You must copy `" + encryptToFilePath + "` file and create " +
			"several backups so as not to lose")
		return
	case "decrypt":
		err := decrypt()
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

func decrypt() error {
	// TODO
	//tHandler := newTimeHandler()
	//tHandler.parseTimeFromWorldAPI()
	return nil
}
