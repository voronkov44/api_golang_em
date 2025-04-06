package enrichment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EnrichmentService interface {
	GetAge(name string) (int, error)
	//GetGender(name string) (string, error)
	//GetNationality(name string) (string, error)
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

//func (e *enrichmentService) GetGender(name string) (string, error) {
//
//}

//func (e *enrichmentService) GetNationality(name string) (string, error) {
//}
