package main

import (
	"flag"
	"log"
)

const (
	encryptFromFilePath = "encrypt_this.txt"
	decryptToFilePath   = "decrypted.txt"
)

func main() {
	flagMode := flag.String("mode", "encrypt", "encrypt / decrypt")

	switch *flagMode {
	default:
		log.Fatalln("unknown mode specified (or not specified?)")
	case "encrypt":
		encrypt()
		return
	case "decrypt":
		decrypt()
		return
	}

	//tHandler := newTimeHandler()
	//tHandler.parseTimeFromWorldAPI()
}

func encrypt() {

}

func decrypt() {

}
