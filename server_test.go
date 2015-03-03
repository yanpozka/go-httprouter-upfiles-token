package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

	fmt.Println("[+] Response: ", respRecorder.Body)
}

//
func TestUploadFile(t *testing.T) {
	mrouter := ConfigRouters()
	respRecorder := httptest.NewRecorder()
	cosa, _ := os.Open(".gitignore")
	req, err := http.NewRequest("POST", "/file", cosa)
	if err != nil {
		t.Fatal("Creating POST '/file' request failed!")
	}

	mrouter.ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusAccepted {
		t.Fatal("Server error: Returned ", respRecorder.Code, " instead of ", http.StatusAccepted)
	}
}
