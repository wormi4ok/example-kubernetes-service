package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/akyoto/cache"

	"github.com/wormi4ok/example-kubernetes-service/opensensemap"
)

type response struct {
	AverageTemperature float64 `json:"averageTemperature"`
}

func temperatureHandler(client *opensensemap.Client, c *cache.Cache, boxIds []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cachedVal, found := c.Get("avg")
		if found {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(response{AverageTemperature: cachedVal.(float64)})
			if err != nil {
				log.Printf("failed to send response: %s", err)
				w.WriteHeader(500)
			}
			return
		}

		// Take data for the last 5 minutes
		fromTime := time.Now().Add(-5 * time.Minute)
		sensors, resp, err := client.BoxesData(r.Context(), boxIds, fromTime, time.Now(), opensensemap.PhenomenonTemperatur)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if resp != nil {
				log.Printf("client got response: %s", resp.Status)
			}
			return
		}
		tt := currentTemperature(sensors)
		avg := averageTemperature(tt)

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response{AverageTemperature: avg})
		if err != nil {
			log.Printf("failed to send response: %s", err)
			w.WriteHeader(500)
		}
		c.Set("avg", avg, 1*time.Minute)
	}
}

// currentTemperature returns a slice of latest values from all sensors given
func currentTemperature(sensors []opensensemap.Sensor) []string {
	latestValues := make(map[string]opensensemap.Sensor)
	for _, sensor := range sensors {
		if val, ok := latestValues[sensor.SensorID]; !ok {
			latestValues[sensor.SensorID] = sensor
		} else if val.CreatedAt < sensor.CreatedAt {
			latestValues[sensor.SensorID] = sensor
		}
	}

	var result []string
	for _, sensorData := range latestValues {
		result = append(result, sensorData.Value)
	}

	return result
}

func averageTemperature(values []string) float64 {
	var sum float64
	var total = float64(len(values))

	for _, value := range values {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil || math.IsNaN(v) {
			log.Printf("ignoring malformed temperature value: %s", value)
			total--
		} else {
			sum += v
		}
	}

	return sum / total
}
