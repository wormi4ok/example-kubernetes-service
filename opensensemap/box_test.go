package opensensemap_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wormi4ok/example-kubernetes-service/opensensemap"
)

func TestClient_BoxesData(t *testing.T) {
	c := opensensemap.NewClient(http.DefaultClient)
	ids := []string{"5cf9874107460b001b828c5b", "5ca4d598cbf9ae001a53051a", "59f8af62356823000fcc460c"}
	fromTime := time.Now().Add(-5 * time.Minute)
	sensors, _, err := c.BoxesData(context.TODO(), ids, fromTime, time.Now(), opensensemap.PhenomenonTemperatur)
	if err != nil {
		t.Fatal(err)
	}
	if len(sensors) == 0 {
		t.Error("Expected data from sensors, got empty response")
	}

	if sensors[0].Value == "" {
		t.Error("Expected Value from sensor, got empty string ")
	}
}

func TestClient_BoxDataMock(t *testing.T) {
	testId := "5cf9874107460b001b828c5b"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, fmt.Sprintf(`[{"sensorId": "%s", "value" : "%s"}]`, "staticID", "42.007"))
	}))
	defer srv.Close()

	c := opensensemap.NewClient(http.DefaultClient)
	c.SetBaseUrl(srv.URL)

	_, _, err := c.BoxesData(context.TODO(), []string{testId}, time.Now().Add(-5*time.Minute), time.Now(), opensensemap.PhenomenonTemperatur)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidation(t *testing.T) {
	type args struct {
		boxIDs   []string
		timeFrom time.Time
		timeTo   time.Time
	}

	tests := []struct {
		name          string
		args          args
		expectedError error
	}{
		{"Basic", args{[]string{"wrongId"}, time.Now().Add(-time.Minute), time.Now()}, opensensemap.ErrInvalidId},
		{"Regex", args{[]string{"cdefghijklmnopqrstuvwxyz"}, time.Now().Add(-time.Minute), time.Now()}, opensensemap.ErrInvalidId},
		{"Date ranges", args{[]string{"5ca4d598cbf9ae001a53051a"}, time.Now(), time.Now().Add(-time.Minute)}, opensensemap.ErrInvalidDatesRange},
	}
	c := opensensemap.NewClient(http.DefaultClient)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, got := c.BoxesData(context.TODO(), tt.args.boxIDs, tt.args.timeFrom, tt.args.timeTo, opensensemap.PhenomenonTemperatur); !errors.Is(got, tt.expectedError) {
				t.Errorf("Expected to get an error: = %v,  got: %v", tt.expectedError, got)
			}
		})
	}
}
