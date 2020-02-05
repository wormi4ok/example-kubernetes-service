package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const myFavouriteTree = "Birch"
const defaultPort = "8080"

func main() {
	port := portFromEnv()

	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "",
	})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}
	view.Register(
		ochttp.ServerRequestCountView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		&view.View{
			Name:        "http_server_response_count_by_status_code",
			Description: "Server response count by status code",
			TagKeys:     []tag.Key{ochttp.StatusCode, ochttp.KeyServerRoute},
			Measure:     ochttp.ServerLatency,
			Aggregation: view.Count(),
		},
	)

	http.Handle("/metrics", ochttp.WithRouteTag(pe, "/metrics"))
	http.Handle("/tree", ochttp.WithRouteTag(treeHandler(myFavouriteTree), "/tree"))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      &ochttp.Handler{},
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Starting service on port %s...\n", port)

	// Running in goroutine so we can shutdown gracefully
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("%v\n", err)
		}
	}()

	handleServerShutdown(srv)
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
			log.Printf("TreeHandler error: %v", err)
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

func handleServerShutdown(srv *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	s := <-signals

	log.Printf("Got %s signal, shutting down server...\n", strings.ToUpper(s.String()))
	// Wait for 5 seconds before shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}
