package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	doc := Document{Name: "Yandry", Size: 260}
	json.NewEncoder(w).Encode(doc)
}

func FileUpload(w http.ResponseWriter, r *http.Request) {

	if r.MultipartForm != nil {
		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				file, _ := fileHeader.Open()
				path := fmt.Sprintf("files/%s", fileHeader.Filename)
				buf, _ := ioutil.ReadAll(file)
				ioutil.WriteFile(path, buf, os.ModePerm)
			}
			fmt.Println("Files founds **************")
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
