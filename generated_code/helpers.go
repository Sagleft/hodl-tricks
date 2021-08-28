package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func httpGET(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to http get: " + err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	return body, nil
}

func readFile(filepath string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.New("failed to read file: " + err.Error())
	}
	return fileBytes, nil
}

func rsaEncrypt(stringToEncrypt string, keyString string) (string, error) {
	// founded at:
	// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes

	// Since the key is in string, we need to convert decode it to bytes
	/*key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", errors.New("failed to decode hex string: " + err.Error())
	}*/
	key := []byte(keyString)
	plaintext := []byte(stringToEncrypt)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("failed to create new cipher from key: " + err.Error())
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("failed to create new gcm from block: " + err.Error())
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("failed to read nonce: " + err.Error())
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data.
	// The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func rsaDecrypt(encryptedString string, keyString string) (string, error) {

	//key, _ := hex.DecodeString(keyString)
	key := []byte(keyString)
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", errors.New("failed to decode hex string to bytes: " + err.Error())
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("failed to create new cipher: " + err.Error())
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("failed to create new gcm from block: " + err.Error())
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("failed to open aes gsm: " + err.Error())
	}

	return fmt.Sprintf("%s", plaintext), nil
}
