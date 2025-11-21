package greet

// HelloResponse represents the response from a hello request
type HelloResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// HelloRequest represents a request to say hello
type HelloRequest struct {
	Name string `json:"name"`
}
