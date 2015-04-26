package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"index", "GET", "/", Index,
	},
	Route{
		"file-upload", "POST", "/file", FileUpload,
	},
}
