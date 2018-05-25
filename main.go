package main

import (
	"engine-manager/lib/kubernetes"
	"engine-manager/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	kube := kubernetes.Init()
	r := mux.NewRouter()
	r.HandleFunc("/sessions", routes.GetPods(kube)).Methods("GET")
	r.HandleFunc("/sessions/{name:[a-zA-Z0-9-_]+}", routes.GetPod(kube)).Methods("GET")
	r.HandleFunc("/sessions/{name:[a-zA-Z0-9-_]+}", routes.CreatePod(kube)).Methods("POST")
	r.HandleFunc("/sessions/{name:[a-zA-Z0-9-_]+}/delete", routes.DeletePod(kube)).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", r))
}
