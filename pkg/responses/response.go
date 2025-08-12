package responses

type DataResponse struct {
	StatusCode int `json:"statusCode"`
	Data       any `json:"data"`
}

type MessageResponse struct {
	StatusCode int `json:"statusCode"`
	Message    any `json:"message"`
}
