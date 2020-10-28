package datamodel

//data model for clWeather
type Response struct {
	Context    string     `json:"@context"`
	ID         string     `json:"id"`
	NWSType    string     `json:"type"`
	Geometry   string     `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	ID         string `json:"@id"`
	NWSType    string `json:"@type"`
	RawMessage string `json:"rawMessage"`
}
