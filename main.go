package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/akyoto/cache"
	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/wormi4ok/example-kubernetes-service/opensensemap"
)

const helloMsg = "Wherever you go, no matter what the weather, always bring your own sunshine.\n"

type Config struct {
	SenseBoxIDs []string `envconfig:"sense_box_ids" required:"true"`
	ListenPort  int      `split_words:"true" default:"8080"`
}

func main() {
	c := new(Config)
	envconfig.MustProcess("", c)

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
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
	)

	httpClient := &http.Client{
		Timeout:   3 * time.Second,
		Transport: &ochttp.Transport{}, // Collect metrics on the HTTP client
	}
	client := opensensemap.NewClient(httpClient)

	http.Handle("/metrics", ochttp.WithRouteTag(pe, "/metrics"))
	http.Handle("/health", ochttp.WithRouteTag(healthHandler(helloMsg), "/hello"))
	http.Handle("/temperature", ochttp.WithRouteTag(temperatureHandler(client, cache.New(time.Minute), c.SenseBoxIDs), "/temperature"))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.ListenPort),
		Handler:      &ochttp.Handler{},
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Starting service on port %d...\n", c.ListenPort)

	// Running in goroutine so we can shutdown gracefully
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("%v\n", err)
		}
	}()

	handleServerShutdown(srv)
}

func healthHandler(helloMsg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, err := io.WriteString(w, helloMsg)
		if err != nil {
			log.Printf("HealthHandler error: %v", err)
			w.WriteHeader(500)
		}
	}
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
