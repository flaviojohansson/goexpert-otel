package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	viaCEPURL  = "https://viacep.com.br/ws/%s/json/"
	weatherURL = "https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no"
)

var weatherAPIKey string // API key da WeatherAPI

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func getTemperatura(localidade string) (float64, float64, float64, error) {

	resp, err := http.Get(fmt.Sprintf(weatherURL, weatherAPIKey, url.QueryEscape(localidade)))
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, 0, fmt.Errorf("error fetching weather data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, err
	}

	var openWeatherResponse WeatherResponse
	if err := json.Unmarshal(body, &openWeatherResponse); err != nil {
		return 0, 0, 0, err
	}

	tempC := openWeatherResponse.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	return tempC, tempF, tempK, nil
}

func getLocalidade(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(viaCEPURL, cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var viaCEPResponse ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResponse); err != nil {
		return "", err
	}

	if viaCEPResponse.Localidade == "" {
		return "", fmt.Errorf("localidade não encontrada")
	}

	return viaCEPResponse.Localidade, nil
}

func climaHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Received request for weather data")
	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 || !isNumeric(cep) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	localidade, err := getLocalidade(cep)
	if err != nil {
		http.Error(w, `{"message": "cannot find zipcode"}`, http.StatusNotFound)
		return
	}

	tempC, tempF, tempK, err := getTemperatura(localidade)
	if err != nil {
		http.Error(w, `{"message": "error fetching weather data"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]float64{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	godotenv.Load()
	weatherAPIKey = os.Getenv("WEATHER_API_KEY")

	http.HandleFunc("/temperatura", climaHandler)
	http.ListenAndServe(":8081", nil)
}
