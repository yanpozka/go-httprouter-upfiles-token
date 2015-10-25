package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

//
type MiddlewareHandler struct {
	Middlewares []CommonMiddleware
	router      HttpRouter
}

//
func (mw *MiddlewareHandler) ServeHTTP(resw http.ResponseWriter, req *http.Request) {

	// custom recovery, it seems deferPanic client recovers first
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[+] Recovering: %+v\nrequest: %+v", r, req)
			debug.PrintStack()
			http.Error(resw, `{"error":"internal"}`, http.StatusInternalServerError)
		}
	}()

	for _, f_mw := range mw.Middlewares {
		if err := f_mw(resw, req); err != nil {
			return
		}
	}

	if mw.router == nil {
		panic("[-] Missing main router.")
	}

	mw.router.ServeHTTP(resw, req) // !!
}

//
type CommonMiddleware func(http.ResponseWriter, *http.Request) error

//
// TODO: remove it. Use https://github.com/julienschmidt/httprouter#basic-authentication
//
func ValidateBasicAuth(u, p string) string {
	if u == "yandry" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

//
func Logger(inner http.Handler, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ruri := r.RequestURI

		//
		inner.ServeHTTP(w, r)

		log.Printf("%s: %s (%s). Time consumed: %s", r.Method, ruri, name, time.Since(start))
	})
}

// Token-accesss protected
func TokenAccessVerification(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if hd, found := r.Header["Authorization"]; !found {
			http.Error(w, "Missing 'Authorization' header", http.StatusUnauthorized)
			return
		} else if len(hd) > 0 && !strings.Contains(hd[0], "Token ") {
			http.Error(w, "Invalid format for 'Authorization' header", http.StatusUnauthorized)
			return
		}

		//
		inner.ServeHTTP(w, r)
	})
}

//
func CommonHeaders(resw http.ResponseWriter, req *http.Request) error {
	header := resw.Header()
	header.Set("Server", "Yandry")
	header.Set("X-Powered-By", "Yandry-Server 0.1")

	// default application/json for all responses
	header.Set("Content-Type", "application/json")

	return nil
}
