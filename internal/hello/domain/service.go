package domain

// HelloService define la interfaz para el servicio de saludo
type HelloService interface {
	SayHello(name string) (string, error)
}

// HelloUseCase define la interfaz para el caso de uso de saludo
type HelloUseCase interface {
	Execute(name string) (*HelloResponse, error)
}

// HelloResponse representa la respuesta del servicio de saludo
type HelloResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}