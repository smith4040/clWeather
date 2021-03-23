package datamodel

import "encoding/json"

// Response is the data model for clWeather
type Response struct {
	Context  []interface{} `json:"@context"`
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID        string `json:"@id"`
		Type      string `json:"@type"`
		Elevation struct {
			Value    int    `json:"value"`
			UnitCode string `json:"unitCode"`
		} `json:"elevation"`
		Station         string `json:"station"`
		Timestamp       string `json:"timestamp"`
		RawMessage      string `json:"rawMessage"`
		TextDescription string `json:"textDescription"`
		Icon            string `json:"icon"`
		PresentWeather  []struct {
			Intensity interface{} `json:"intensity"`
			Modifier  interface{} `json:"modifier"`
			Weather   string      `json:"weather"`
			RawString string      `json:"rawString"`
		} `json:"presentWeather"`
		Temperature struct {
			Value          JSONInt `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"temperature"`
		Dewpoint struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"dewpoint"`
		WindDirection struct {
			Value          int    `json:"value"`
			UnitCode       string `json:"unitCode"`
			QualityControl string `json:"qualityControl"`
		} `json:"windDirection"`
		WindSpeed struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"windSpeed"`
		WindGust struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl string      `json:"qualityControl"`
		} `json:"windGust"`
		BarometricPressure struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"barometricPressure"`
		SeaLevelPressure struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl string      `json:"qualityControl"`
		} `json:"seaLevelPressure"`
		Visibility struct {
			Value          int    `json:"value"`
			UnitCode       string `json:"unitCode"`
			QualityControl string `json:"qualityControl"`
		} `json:"visibility"`
		MaxTemperatureLast24Hours struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl interface{} `json:"qualityControl"`
		} `json:"maxTemperatureLast24Hours"`
		MinTemperatureLast24Hours struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl interface{} `json:"qualityControl"`
		} `json:"minTemperatureLast24Hours"`
		PrecipitationLastHour struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl string      `json:"qualityControl"`
		} `json:"precipitationLastHour"`
		PrecipitationLast3Hours struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl string      `json:"qualityControl"`
		} `json:"precipitationLast3Hours"`
		PrecipitationLast6Hours struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl string      `json:"qualityControl"`
		} `json:"precipitationLast6Hours"`
		RelativeHumidity struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"relativeHumidity"`
		WindChill struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"windChill"`
		HeatIndex struct {
			Value          interface{} `json:"value"`
			UnitCode       string      `json:"unitCode"`
			QualityControl string      `json:"qualityControl"`
		} `json:"heatIndex"`
		CloudLayers []struct {
			Base struct {
				Value    int    `json:"value"`
				UnitCode string `json:"unitCode"`
			} `json:"base"`
			Amount string `json:"amount"`
		} `json:"cloudLayers"`
	} `json:"properties"`
}

// JSONInt is a special struct for handling null vs 0 deg temps
type JSONInt struct {
	Value float64
	Valid bool
	Set   bool
}

// UnmarshalJSON is a method for handling null vs 0 deg values from API
func (i *JSONInt) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true

	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}

	// The key isn't set to null
	var temp float64
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
