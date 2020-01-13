package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const myFavouriteTree = "Birch"

func main() {
	http.HandleFunc("/tree", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := io.WriteString(w, fmt.Sprintf("{\"myFavouriteTree\":\"%s\"}\n", myFavouriteTree))
		if err != nil {
			w.WriteHeader(500)
		}
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", "8000"),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Print("Starting service on port 8000")

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("%v\n", err)
	}
}
