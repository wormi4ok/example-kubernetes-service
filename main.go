package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const myFavouriteTree = "Birch"
const defaultPort = "8080"

func main() {
	port := portFromEnv()

	http.HandleFunc("/tree", treeHandler(myFavouriteTree))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Starting service on port %s...\n", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("%v\n", err)
	}
}

type response struct {
	MyFavouriteTree string `json:"myFavouriteTree"`
}

func treeHandler(favouriteTree string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(404)
			return
		}
		w.Header().Set("Content-Type", "application/json") // RFC-4627
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(response{MyFavouriteTree: favouriteTree})
		if err != nil {
			w.WriteHeader(500)
		}
	}
}

func portFromEnv() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		log.Printf("PORT env var is not specified, using default: %s", port)
	}
	return port
}
