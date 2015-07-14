package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

const PORT = ":8080"

//
func main() {

	fmt.Println("[+] Init webserver", PORT)

	log.Fatal(http.ListenAndServe(PORT, ConfigRouters()))
}

//
// https://github.com/corylanou/tns-restful-json-api/blob/master/v9/logger.go
//
func Logger(inner http.Handler, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

//
func ConfigRouters() *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		var handler http.Handler = Logger(route.HandlerFunc, route.Name)

		router.Handler(route.Method, route.Pattern, handler)
	}

	return router
}
