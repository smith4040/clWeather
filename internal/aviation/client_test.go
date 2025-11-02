package aviation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchMETAR_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"METAR":[{"raw_text":"KJFK 011200Z ...","station_id":"KJFK","fltCat":"VFR"}]}}`))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	oldURL := baseURL
	baseURL = server.URL // Mock base
	defer func() { baseURL = oldURL }()

	data, err := FetchMETAR(ctx, "KJFK")
	require.NoError(t, err)
	assert.Contains(t, string(data), "KJFK")
}

func TestParseMETAR(t *testing.T) {
	jsonData := []byte(`{
        "data": {
            "METAR": [{
                "raw_text": "KJFK 011200Z 27010G15KT 10SM FEW030 BKN050 15/10 A2992",
                "station_id": "KJFK",
                "observation_time": "2025-11-01T12:00:00Z",
                "temperature": "15",
                "dewpoint": "10",
                "wind": {"direction": "270", "speed": "10", "gust": "15"},
                "visibility": {"statute_miles": "10"},
                "clouds": [{"type": "FEW", "base": 3000}, {"type": "BKN", "base": 5000}],
                "weather": [{"type": "CLR"}],
                "barometer": {"value": 29.92},
                "elevation": {"value": 13},
                "fltCat": "VFR"
            }]
        }
    }`)

	r, err := ParseMETAR(jsonData)
	require.NoError(t, err)
	assert.Equal(t, "KJFK", r.StationID)
	assert.Equal(t, "VFR", r.FltCat)
	assert.Equal(t, "15°C", r.Temperature)
	assert.Len(t, r.Clouds, 2)
	assert.Equal(t, "FEW", r.Clouds[0].Type)
}

func TestParseTAF(t *testing.T) {
	jsonData := []byte(`{
        "data": {
            "TAF": [{
                "raw_text": "KJFK 011130Z 0112/0212 ...",
                "station_id": "KJFK",
                "issue_time": "2025-11-01T11:30:00Z",
                "valid_time_from": "2025-11-01T12:00:00Z",
                "valid_time_to": "2025-11-02T12:00:00Z",
                "forecast": [{
                    "type": "forecast",
                    "start_time": "2025-11-01T12:00:00Z",
                    "end_time": "2025-11-02T12:00:00Z",
                    "wind": {"direction": "270", "speed": "8"},
                    "visibility": {"statute_miles": "6+"},
                    "clouds": [{"type": "FEW", "base": 4000}]
                }]
            }]
        }
    }`)

	r, err := ParseTAF(jsonData)
	require.NoError(t, err)
	assert.Equal(t, "KJFK", r.StationID)
	assert.Len(t, r.Forecast, 1)
	assert.Equal(t, "forecast", r.Forecast[0].Type)
}
