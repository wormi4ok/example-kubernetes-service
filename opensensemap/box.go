package opensensemap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const baseUrl = "https://api.opensensemap.org"

// Client to use OpenSenseMap API
type Client struct {
	baseUrl    string
	baseClient *http.Client
}

// Phenomenon represents a weather phenomenon in the API
type Phenomenon string

const PhenomenonTemperatur Phenomenon = "Temperatur"

type Sensor struct {
	SensorID  string  `json:"sensorId"`
	Value     string  `json:"value"`
	CreatedAt string  `json:"createdAt"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
}

// NewClient creates a new OpenSenseMap API client
func NewClient(client *http.Client) *Client {
	return &Client{
		baseUrl:    baseUrl,
		baseClient: client,
	}
}

// SetBaseUrl allows user to override the API hostname specified in `baseUrl`
func (c *Client) SetBaseUrl(url string) {
	c.baseUrl = url
}

var ErrInvalidDatesRange = errors.New("to-date cannot be before from-date")

// BoxesData allows to download data of a given phenomenon from multiple senseBoxes
// This method obtains and parses data in JSON format from sensors.
// There is also a CSV format available to download a analytical data
func (c *Client) BoxesData(
	ctx context.Context,
	senseBoxIds []string,
	fromDate time.Time,
	toDate time.Time,
	phenomenon Phenomenon,
) ([]Sensor, *http.Response, error) {
	for _, id := range senseBoxIds {
		if err := validateId(id); err != nil {
			return nil, nil, err
		}
	}
	if toDate.Before(fromDate) {
		return nil, nil, ErrInvalidDatesRange
	}

	query := url.Values{}
	query.Add("boxId", strings.Join(senseBoxIds, ","))
	query.Add("from-date", fromDate.UTC().Format("2006-01-02T03:04:05Z"))
	query.Add("to-date", toDate.UTC().Format("2006-01-02T03:04:05Z"))
	query.Add("phenomenon", string(phenomenon))
	query.Add("download", "false")
	query.Add("format", "json")

	resp, err := c.doGetRequest(ctx, fmt.Sprintf(c.baseUrl+"/boxes/data?%s", query.Encode()))
	if err != nil {
		return nil, resp, err
	}

	var sensors []Sensor
	err = json.NewDecoder(resp.Body).Decode(&sensors)
	if err != nil {
		return nil, resp, err
	}

	return sensors, resp, nil
}

var ErrResponseNotOK = errors.New("unexpected response code")

func (c *Client) doGetRequest(ctx context.Context, urlEncoded string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEncoded, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("%w:%s", ErrResponseNotOK, resp.Status)
	}
	return resp, nil
}

var ErrInvalidId = errors.New("incorrect ID format")

func validateId(id string) error {
	matched, err := regexp.MatchString(`^[\da-f]{24}`, id)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("validate Id %v: %w", id, ErrInvalidId)
	}

	return nil
}
