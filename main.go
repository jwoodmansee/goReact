package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{searchParams}", index).Methods("GET")

	staticFileDirectory := http.Dir("./dist")
	staticFileHandler := http.StripPrefix("/dist/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/dist/").Handler(staticFileHandler).Methods("GET")
	return r
}

func main() {
	fmt.Println("Server is running")
	r := newRouter()

	log.Fatal(http.ListenAndServe(getPort(), handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}
