package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type FindWeatherHandler struct {
	WeatherAPIKey string
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type Weather struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type Temperatures struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewFindWeatherHandler(weatherAPIKey string) *FindWeatherHandler {
	return &FindWeatherHandler{
		WeatherAPIKey: weatherAPIKey,
	}
}

func (h *FindWeatherHandler) FindWeather(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		http.Error(w, "CEP is required", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valid, err := VerifyValidCEP(cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	cep, err := GetCEPInfo(cepParam)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			http.Error(w, err.Error(), http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, "Error getting CEP info: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weather, err := GetWeatherInfo(cep.Localidade, h.WeatherAPIKey)
	if err != nil {
		http.Error(w, "Error getting weather info: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temperatures := Temperatures{
		TempC: weather.Current.TempC,
		TempF: weather.Current.TempF,
		TempK: weather.Current.TempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(temperatures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func VerifyValidCEP(cep string) (bool, error) {
	valid, err := regexp.MatchString("\\d{5}-*\\d{3}", cep)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func GetCEPInfo(cep string) (*ViaCEP, error) {
	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "erro") {
		return nil, errors.New("can not find zipcode")
	}

	var cepData ViaCEP

	err = json.Unmarshal(body, &cepData)
	if err != nil {
		return nil, err
	}

	return &cepData, nil
}

func GetWeatherInfo(location string, apiKey string) (*Weather, error) {
	requestUrl := "http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + url.QueryEscape(location)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "erro") {
		var errorResponse struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(errorResponse.Error.Message)
	}

	var weather Weather

	err = json.Unmarshal(body, &weather)

	if err != nil {
		return nil, err
	}

	return &weather, nil
}
