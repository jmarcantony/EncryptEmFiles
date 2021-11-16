package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	dataFile = "data.json"
	port     = ":8080"
)

type result map[string]string

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	http.HandleFunc("/get", get)
	http.ListenAndServe(port, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK!")
}

func add(w http.ResponseWriter, r *http.Request) {
	var d result
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Fatal(err)
	}
	appendJson(d)
}

func get(w http.ResponseWriter, r *http.Request) {
	id, _ := r.URL.Query()["id"]
	f, err := ioutil.ReadFile(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	var data result
	_ = json.Unmarshal(f, &data)
	fmt.Fprintf(w, data[id[0]])
}

func appendJson(d result) {
	f, err := ioutil.ReadFile(dataFile)
	if err != nil {
		data, _ := json.MarshalIndent(d, "", "    ")
		_ = ioutil.WriteFile(dataFile, data, 0644)
		return
	}
	var init result
	_ = json.Unmarshal(f, &init)
	for key, value := range d {
		init[key] = value
	}
	data, _ := json.MarshalIndent(init, "", "    ")
	_ = ioutil.WriteFile(dataFile, data, 0644)
}
