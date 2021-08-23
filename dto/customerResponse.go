package dto

type CustomerResponse struct {
	Id      string `json:"uuid"`
	Name    string `json:"name"`
	City    string `json:"city"`
	ZipCode string `json:"postalCode"`
	DOB     string `json:"dob"`
	Status  string `json:"state"`
}
