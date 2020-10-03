package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/akyoto/cache"

	"github.com/wormi4ok/example-kubernetes-service/opensensemap"
)

func Test_temperatureHandler(t *testing.T) {
	ids := []string{"5cf9874107460b001b828c5b", "5ca4d598cbf9ae001a53051a", "59f8af62356823000fcc460c"}
	c := opensensemap.NewClient(http.DefaultClient)
	srv := httptest.NewServer(temperatureHandler(c, cache.New(time.Minute), ids))
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, srv.URL+"/temperature", nil)
	req.RequestURI = "" // Request.RequestURI can't be set in client requests
	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that HTTP status code is 200
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code = %d, actual = %d", http.StatusOK, res.StatusCode)
	}

	// Check that the response format is correct
	var got response
	err = json.NewDecoder(res.Body).Decode(&got)
	if err != nil {
		t.Errorf("parsing response: %s", err)
	}
}

func Test_currentTemperature(t *testing.T) {
	tests := []struct {
		name    string
		sensors []opensensemap.Sensor
		want    []string
	}{
		{"Return latest for sensor",
			[]opensensemap.Sensor{{
				SensorID:  "sensor1",
				Value:     "10",
				CreatedAt: "2020-10-3T15:53:09.829Z",
			}, {
				SensorID:  "sensor1",
				Value:     "5",
				CreatedAt: "2020-10-3T15:55:09.829Z",
			}}, []string{"5"},
		},
		{"Return value for each sensor",
			[]opensensemap.Sensor{{
				SensorID:  "sensor2",
				Value:     "20",
				CreatedAt: "2020-10-3T15:57:09.829Z",
			}, {
				SensorID:  "sensor3",
				Value:     "15",
				CreatedAt: "2020-10-3T15:58:09.829Z",
			}}, []string{"15", "20"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := currentTemperature(tt.sensors)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("currentTemperature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_averageTemperature(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   float64
	}{
		{"Happy path", []string{"7.5", "21.5"}, 14.5},
		{"Special number", []string{"21", "NaN"}, 21},
		{"Invalid input", []string{"1", "cold", "6"}, 3.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := averageTemperature(tt.values); got != tt.want {
				t.Errorf("averageTemperature() = %v, want %v", got, tt.want)
			}
		})
	}
}
