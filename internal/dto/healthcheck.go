package dto

type HealthcheckResponse struct {
	Status string   `json:"status"`
	Errors []string `json:"errors"`
}
