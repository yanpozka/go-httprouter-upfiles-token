package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Index(w http.ResponseWriter, r *http.Request) {
	doc := Document{Name: "Yandry", Size: 260}
	json.NewEncoder(w).Encode(doc)
}

func FileUpload(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("docfile")

	if err != nil {
		fmt.Println("[-] Error in r.FormFile ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[-] Error in r.FormFile ", err)
		return
	}
	defer file.Close()

	out, err := os.Create("uploaded-" + header.Filename)
	if err != nil {
		fmt.Println("[-] Unable to create the file for writing. Check your write access privilege.", err)
		fmt.Fprintf(w, "[-] Unable to create the file for writing. Check your write access privilege.", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Fail on io.Copy", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println("[+] File uploaded successfully: ")
	fmt.Println("uploaded-" + header.Filename)
}
