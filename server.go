package main 

import (
    "fmt"
    "net/http"
    "log"

    "github.com/bidder/db"
    "github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Setting up server")
}

func main(){
    defer db.MgoSession.Close()
    fmt.Println("Starting go service...")
    
    r := mux.NewRouter()
    r.HandleFunc("/", indexHandler)
    if err := http.ListenAndServe(":3000", r); err != nil {
    	log.Fatal(err)
    }
}