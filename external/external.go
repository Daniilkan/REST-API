package external

// Package external provides utility functions to fetch external data from APIs.
// It includes functions to get age, gender, and nationality based on a person's name.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GetAge fetches the estimated age for a given name from an external API.
// @Summary Get estimated age
// @Description GetAge fetches the estimated age for a given name using the Agify API.
// @Tags external
// @Param name query string true "Name to estimate age for"
// @Success 200 {integer} int "Estimated age"
// @Failure 500 {string} string "Failed to fetch age"
// @Router /external/age [get]
func GetAge(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io?name=%s", name))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	age, ok := result["age"].(float64)
	if !ok {
		return 0, errors.New("age is not a number")
	}
	return int(age), nil
}

// GetGender fetches the gender for a given name from an external API.
// @Summary Get gender
// @Description GetGender fetches the gender for a given name using the Genderize API.
// @Tags external
// @Param name query string true "Name to fetch gender for"
// @Success 200 {string} string "Gender"
// @Failure 500 {string} string "Failed to fetch gender"
// @Router /external/gender [get]
func GetGender(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	gender, ok := result["gender"].(string)
	if !ok {
		return "", errors.New("failed to fetch gender")
	}
	return gender, nil
}

type NationalityResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

// GetNationality fetches the nationality for a given name from an external API.
// @Summary Get nationality
// @Description GetNationality fetches the nationality for a given name using the Nationalize API.
// @Tags external
// @Param name query string true "Name to fetch nationality for"
// @Success 200 {string} string "Nationality"
// @Failure 500 {string} string "Failed to fetch nationality"
// @Router /external/nationality [get]
func GetNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result NationalityResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Country) == 0 {
		return "", errors.New("no countries found in response")
	}

	return result.Country[0].CountryID, nil
}
