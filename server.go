package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//
func main() {

	fmt.Println("[+] Init webserver")

	log.Fatal(http.ListenAndServe(":8080", ConfigRouters()))
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
func ConfigRouters() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Methods).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
