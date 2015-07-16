package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"index", "GET", "/", Index,
	},
	Route{
		"index1", "GET", "/index", Index,
	},
	Route{
		"file-upload", "POST", "/file", FileUpload,
	},
	Route{
		"generate-token", "GET", "/access-token", GenerateSecurityToken,
	},
}
