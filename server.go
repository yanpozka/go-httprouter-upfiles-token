package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

const PORT = ":8080"

// Interface for any router: httprouter, gorilla-mux, etc.
type HttpRouter interface {
	http.Handler
}

//
type MiddlewareHandler struct {
	Middlewares []CommonMiddleware
	router      HttpRouter
}

//
func (mw *MiddlewareHandler) ServeHTTP(resw http.ResponseWriter, req *http.Request) {

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
func main() {
	log.Printf("[+] Starting server in %s\n", PORT)

	log.Fatal(http.ListenAndServe(PORT, newApp()))
}

//
func newApp() *MiddlewareHandler {

	// Logger, Common Headers middlewares
	mdws := []CommonMiddleware{CommonHeaders}

	mwhanderl := &MiddlewareHandler{Middlewares: mdws, router: ConfigRouters()}

	return mwhanderl
}

//
func ConfigRouters() *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		var handler http.Handler = Logger(TokenAccessVerification(route.HandlerFunc), route.Name)

		router.Handler(route.Method, route.Path, handler)
		log.Printf("[+] Registred endpoint %s: %s (%s)", route.Method, route.Path, route.Name)
	}
	// access-token endpoint
	router.Handler(rat.Method, rat.Path, Logger(rat.HandlerFunc, rat.Name))
	log.Printf("[+] Registred endpoint %s: %s (%s)", rat.Method, rat.Path, rat.Name)

	return router
}
