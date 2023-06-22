package models

type ApiShortenRequest struct {
	Url string `json:"url"`
}

type ApiShortenResponse struct {
	Result string `json:"result"`
}
