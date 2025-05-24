package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
	switch strings.ToLower(gender) {
	case "male":
		return "m", nil
	case "female":
		return "f", nil
	default:
		return "", errors.New("unknown gender")
	}
}

type NationalityResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

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
