package enrichment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EnrichmentService interface {
	GetAge(name string) (int, error)
	GetGender(name string) (string, error)
	GetNationality(name string) (string, error)
}

type enrichmentService struct{}

func NewEnrichmentService() EnrichmentService {
	return &enrichmentService{}
}

func (e *enrichmentService) GetAge(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Age int `json:"age"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Age, err
}

func (e *enrichmentService) GetGender(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "error request", err
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Gender, err
}

func (e *enrichmentService) GetNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "error request", err
	}
	defer resp.Body.Close()

	var result struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil || len(result.Country) == 0 {
		return "error request", err
	}
	return result.Country[0].CountryID, nil
}
