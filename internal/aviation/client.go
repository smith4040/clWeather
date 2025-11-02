package aviation

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://aviationweather.gov/api/data"

var client = &http.Client{
	Timeout: 8 * time.Second,
	Transport: &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: true,
	},
}

func FetchMETAR(ctx context.Context, station string) ([]byte, error) {
	return fetch(ctx, "metar", station)
}

func FetchTAF(ctx context.Context, station string) ([]byte, error) {
	return fetch(ctx, "taf", station)
}

func fetch(ctx context.Context, endpoint, station string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "aviationweather.gov",
		Path:   fmt.Sprintf("/api/data/%s", endpoint),
	}
	q := u.Query()
	q.Set("ids", station)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "clWeather-Aviation/1.0 (+https://github.com/smith4040/clWeather)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded (HTTP 429); try again in 1 min")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("no data for station %s", station)
	}

	return body, nil
}
