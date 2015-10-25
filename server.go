package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	auth "github.com/abbot/go-http-auth"
	"github.com/deferpanic/deferclient/deferstats"
	"github.com/julienschmidt/httprouter"
)

const PORT = ":8080"

// Interface for any router: httprouter, gorilla-mux, etc.
type HttpRouter interface {
	http.Handler
}

//
func main() {
	log.Printf("[+] Starting server in %s\n", PORT)

	log.Fatal(http.ListenAndServe(PORT, newApp()))
}

//
// DeferPanic examples
//
func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("there is no need to PANIC")
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this request is FAST")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Fprintf(w, "this request is SLOW")
}

//
func newApp() *MiddlewareHandler {
	// Logger, Common Headers middlewares
	mdws := []CommonMiddleware{CommonHeaders}

	router := ConfigRouters()

	if os.Getenv("IS_TESTING") == "" {
		dps := deferstats.NewClient(os.Getenv("DEFERPANIC_API_KEY"))

		go dps.CaptureStats()

		router.Handler("GET", "/fast", dps.HTTPHandlerFunc(fastHandler))
		router.Handler("GET", "/slow", dps.HTTPHandlerFunc(slowHandler))
		router.Handler("GET", "/panic", dps.HTTPHandlerFunc(panicHandler))
	}

	mwhanderl := &MiddlewareHandler{Middlewares: mdws, router: router}

	return mwhanderl
}

//
func ConfigRouters() *httprouter.Router {
	router := httprouter.New()

	ar := auth.NewBasicAuthenticator("localhost", ValidateBasicAuth)

	for _, route := range routes {
		var handler http.Handler = Logger(TokenAccessVerification(route.HandlerFunc), route.Name)

		router.Handler(route.Method, route.Path, handler)
		log.Printf("[+] Registred endpoint %s: %s (%s)", route.Method, route.Path, route.Name)
	}

	// access-token endpoint
	router.Handler(rat.Method, rat.Path, Logger(ar.Wrap(GenerateSecurityToken), rat.Name))

	log.Printf("[+] Registred endpoint %s: %s (%s)", rat.Method, rat.Path, rat.Name)

	return router
}
