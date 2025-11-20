package recovery

type ValidateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}