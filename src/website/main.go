package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Daily struct {
	TemperatureMax []float64 `json:"temperature_2m_max"`
	TemperatureMin []float64 `json:"temperature_2m_min"`
	Sunrise        []string  `json:"sunrise"`
	Sunset         []string  `json:"sunset"`
}

type WeatherResponse struct {
	Daily Daily `json:"daily"`
}

const url = "https://api.open-meteo.com/v1/forecast?latitude=54.6892&longitude=25.2798&daily=temperature_2m_max,temperature_2m_min,sunrise,sunset,precipitation_hours&hourly=temperature_2m,cloud_cover,apparent_temperature,precipitation_probability,precipitation&timezone=auto&forecast_days=1"

func handler(w http.ResponseWriter, r *http.Request) {
	response, error := http.Get(url)
	if error != nil {
		fmt.Println("Error fetching weather:", error)
		return
	}
	defer response.Body.Close()

	body, error := io.ReadAll(response.Body)
	if error != nil {
		fmt.Println("Error reading response:", error)
		return
	}

	var weather WeatherResponse
	json.Unmarshal(body, &weather)

	var prettyJSON bytes.Buffer
	error = json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
		return
	}
	fmt.Fprintln(w, prettyJSON.String())
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("moonbase listening on :8080")
	http.ListenAndServe(":8080", nil)
}
