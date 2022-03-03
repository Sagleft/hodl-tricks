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
	timeLayoutTt        = "2006-01-02 15:04:05.999999999 -0700 MST"
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
		log.Println("(!!) The next run of the utility will overwrite the file." +
			"Therefore, if you have changed the data file, be careful")
		return
	case "decrypt":
		err := decrypt()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("YES! Time's up! File decrypted successfuly!")
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
	// read file
	fileBytes, err := readFile(encryptToFilePath)
	if err != nil {
		return err
	}

	// decrypt file
	decryptedBytes, err := rsaDecrypt(fileBytes, encryptionKey)
	if err != nil {
		return err
	}

	// json decode
	data := dataContainer{}
	err = json.Unmarshal(decryptedBytes, &data)
	if err != nil {
		return errors.New("failed to unmarshal data container json: " +
			err.Error())
	}

	// parse unlockOn timestamp
	unlockOnTimeParsed, err := time.Parse(timeLayoutTt, data.UnlockOn)
	if err != nil {
		return errors.New("failed to parse time: " + err.Error())
	}

	// get timestamp from external API
	tHandler := newTimeHandler()
	timeFromAPI, err := tHandler.getCurrentTime()
	if err != nil {
		return err
	}

	durationDelta := unlockOnTimeParsed.Sub(*timeFromAPI)
	if durationDelta > 0 {
		return errors.New("Time is not up yet. Come back through " + durationDelta.String())
	}

	// time's up! save decoded file
	err = saveToFile(decryptToFilePath, data.Data)
	if err != nil {
		return err
	}
	return nil
}
