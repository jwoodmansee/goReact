package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func GetAmpSpecs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	//Change the varible to something better
	stringSlice := strings.Split(r.URL.Path, "/")
	//ampInfo := r.URL.Query().Get("params")

	if len(stringSlice[2]) > 0 {
		fmt.Println("You did it...kind of!")
	}

}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/search/{params}", GetAmpSpecs).Methods("GET")

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
