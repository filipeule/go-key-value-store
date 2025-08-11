package main

import "github.com/gorilla/mux"

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", helloMuxHandler)

	r.HandleFunc("/v1/key/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", keyValueDeleteHandler).Methods("DELETE")

	return r
}