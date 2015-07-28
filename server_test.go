package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//
var (
	mrouter    HttpRouter
	token_hash string
)

//
func init() {
	mrouter = newApp()

	resp, err := _getToken()
	if err != nil {
		panic(fmt.Sprintf("Imposible create Token. %s", err))
		return
	}

	var tr = &struct {
		Token string `json:"token"`
	}{}

	if err := json.Unmarshal(resp, tr); err != nil {
		panic(fmt.Sprintf("[-] Error trying to decode token response. %v", err))
	}
	token_hash = tr.Token
}

//
func TestAccessToken(t *testing.T) {
	token, err := _getToken()
	if err != nil {
		t.Fatalf("Error getting Access-Token. %v\n", err)
	}

	if !strings.Contains(string(token), `"token":`) {
		t.Fatalf("[-] Has to Contains token field")
	}
}

//
func _getToken() ([]byte, error) {
	respRec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/access-token", nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Creating GET '/' request failed!")
	}
	req.SetBasicAuth("yandry", "hello") // !!

	mrouter.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		return nil, fmt.Errorf("Returned code %d instead of %d. Server says: %s",
			respRec.Code, http.StatusOK, respRec.Body.String())
	}
	return respRec.Body.Bytes(), nil
}

//
func TestEchosContent(t *testing.T) {
	respRec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Creating GET '/' request failed!")
	}
	req.Header.Add("Authorization", "Token "+token_hash)

	mrouter.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusOK)
	}

	header := respRec.HeaderMap
	expectedHeaders := []string{"Content-Type", "X-Powered-By", "Server"}

	for _, chr := range expectedHeaders {
		if _, found := header[chr]; !found {
			t.Logf("[-] Not found expected header '%s'. Header: %v", chr, header)
		}
	}
	fmt.Println("[+] Response: ", respRec.Body)
}

//
func TestUploadFile(t *testing.T) {
	respRecorder := httptest.NewRecorder()
	file_to_upload, errf := os.Open(".gitignore")

	if errf != nil {
		t.Fatal("[-] Fail to open file to upload.", errf)
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

	if err := writer.Close(); err != nil {
		t.Fatal("[-] Fail to writer.Close()", err)
	}

	req, err := http.NewRequest("POST", "/file", body)
	if err != nil {
		t.Fatal("[-] Creating POST '/file' request failed!")
	}
	req.Header.Add("Content-Type", writer.FormDataContentType()) // BLOODY LINE OF CODE
	req.Header.Add("Authorization", "Token "+token_hash)

	mrouter.ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusOK {
		t.Fatal("[-] Server error: Returned [", respRecorder.Code, "] instead of [", http.StatusOK, "]")
	}

	t.Log("Code :", respRecorder.Code)
}
