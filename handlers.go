package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO: add Basic Auth to this endpoint
func GenerateSecurityToken(w http.ResponseWriter, r *http.Request) {
	buffrand := make([]byte, 65)

	if _, err := rand.Read(buffrand); err != nil {
		log.Println("[-] Error trying to generate random string.", err)
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}

	// TODO: create a separated function to create hash
	hasher512 := sha512.New()
	var rs string = fmt.Sprintf("%s-%d-%s", r.RemoteAddr, time.Now().UnixNano(), string(buffrand))
	data := hasher512.Sum([]byte(rs))

	fmt.Fprintf(w, `{"token": "%x"}`, data)
}

//
func Index(w http.ResponseWriter, r *http.Request) {
	doc := Document{Name: "Yandry", Size: 260}
	json.NewEncoder(w).Encode(doc)
}

//
func FileUpload(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("docfile")

	if err != nil {
		log.Println("[-] Error in r.FormFile ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{'error': %s}", err)
		return
	}
	defer file.Close()

	out, err := os.Create("uploaded-" + header.Filename)
	if err != nil {
		log.Println("[-] Unable to create the file for writing. Check your write access privilege.", err)
		fmt.Fprintf(w, "[-] Unable to create the file for writing. Check your write access privilege.", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("[-] Error copying file.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("[+] File uploaded successfully: uploaded-", header.Filename)
}
