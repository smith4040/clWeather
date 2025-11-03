package aviation

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func FetchMETAR(ctx context.Context, station string) ([]byte, error) {
	return fetch(ctx, "metar", station)
}

func FetchTAF(ctx context.Context, station string) ([]byte, error) {
	return fetch(ctx, "taf", station)
}

func fetch(ctx context.Context, endpoint, station string) ([]byte, error) {
	u := url.URL{Scheme: "https", Host: "aviationweather.gov", Path: "/api/data/" + endpoint}
	q := u.Query()
	q.Set("ids", station)
	q.Set("format", "json")
	q.Set("mostRecent", "true")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	req.Header.Set("User-Agent", "clWeather/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limit (429)")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
