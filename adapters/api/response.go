package api

type Response struct {
	Message interface{} `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}
