package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

//
func TestEchosContent(t *testing.T) {
	mrouter := ConfigRouters()
	respRecorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Creating GET '/' request failed!")
	}

	mrouter.ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", respRecorder.Code, " instead of ", http.StatusOK)
	}
	// fmt.Println("[+] Response: ", respRecorder.Body)
}

//
func TestUploadFile(t *testing.T) {
	mrouter := ConfigRouters()
	respRecorder := httptest.NewRecorder()
	file_to_upload, errf := os.Open(".gitignore")

	if errf != nil {
		t.Fatal("[-] Fail to open './document.go'", errf)
	}
	defer file_to_upload.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("docfile", filepath.Base("/tmp/gitignore"))
	if err != nil {
		t.Fatal("[-] Fail to writer.CreateFormFile", err)
	}
	_, err = io.Copy(part, file_to_upload)
	if err != nil {
		t.Fatal("[-] Fail to io.Copy(part, file_to_upload)", err)
	}
	// if we want more parameters to send
	// for key, val := range params { _ = writer.WriteField(key, val) }*/

	err = writer.Close()
	if err != nil {
		t.Fatal("[-] Fail to writer.Close()", err)
	}

	req, err := http.NewRequest("POST", "/file", body)
	if err != nil {
		t.Fatal("[-] Creating POST '/file' request failed!")
	}
	req.Header.Add("Content-Type", writer.FormDataContentType()) // BLOODY LINE OF CODE

	mrouter.ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusOK {
		t.Fatal("[-] Server error: Returned [", respRecorder.Code, "] instead of [", http.StatusOK, "]")
	}

	fmt.Println("Code :", respRecorder.Code)
}
