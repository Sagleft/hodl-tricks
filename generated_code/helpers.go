package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	// to check that the IP falls within the range of the private
	// reference: https://en.wikipedia.org/wiki/Private_network
	// from subnet -> to subnet
	ipSubnets = map[int][]string{
		0: []string{"10.0.0.0", "10.255.255.255"},     // single class A network
		1: []string{"100.64.0.0", "100.127.255.255"},  // additional network
		2: []string{"172.16.0.0", "172.31.255.255"},   // 16 contiguous class B networks
		3: []string{"192.168.0.0", "192.168.255.255"}, // 256 contiguous class C networks
		4: []string{"127.0.0.0", "127.255.255.255"},   // localhost
	}
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

func rsaEncrypt(dataToEncrypt []byte, keyString string) ([]byte, error) {
	// founded at:
	// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes

	// Since the key is in string, we need to convert decode it to bytes
	/*key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", errors.New("failed to decode hex string: " + err.Error())
	}*/
	key := []byte(keyString)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create new cipher from key: " + err.Error())
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create new gcm from block: " + err.Error())
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to read nonce: " + err.Error())
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data.
	// The first nonce argument in Seal is the prefix.
	return aesGCM.Seal(nonce, nonce, dataToEncrypt, nil), nil
}

func rsaDecrypt(encrypted []byte, keyString string) ([]byte, error) {

	//key, _ := hex.DecodeString(keyString)
	key := []byte(keyString)
	/*enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", errors.New("failed to decode hex string to bytes: " + err.Error())
	}*/

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create new cipher: " + err.Error())
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create new gcm from block: " + err.Error())
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// Decrypt the data
	result, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			return nil, errors.New("Incorrect decryption key. Are you using the same version of the utility as before?")
		}
		return nil, errors.New("failed to open aes gsm: " + err.Error())
	}

	return result, nil
}

func saveToFile(filepath string, content []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		return errors.New("failed to create file: " + err.Error())
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return errors.New("failed to write string to file: " + err.Error())
	}
	return nil
}

// check host IP is local
func checkHostIP(host string) (bool, error) {
	// get host IPs
	ips, err := net.LookupIP(host)
	if err != nil {
		return false, errors.New("Could not get IPs for host: " + err.Error())
	}
	if len(ips) == 0 {
		return false, errors.New("host IPs not found")
	}

	// get host first IP
	hostIP := ips[0]

	// check IP subnet
	for _, subnet := range ipSubnets {
		ipFrom := net.ParseIP(subnet[0])
		ipTo := net.ParseIP(subnet[1])

		// check IP format
		if hostIP.To4() == nil {
			return false, errors.New("failed to get IPv4 for host: " + host)
		}

		if bytes.Compare(hostIP, ipFrom) >= 0 && bytes.Compare(hostIP, ipTo) <= 0 {
			return true, nil
		}
	}
	return false, nil
}
