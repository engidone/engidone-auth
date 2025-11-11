package endpoints

import (
	"context"

	"engidone-auth/internal/hello/domain"
	"github.com/go-kit/kit/endpoint"
)

// HelloRequest represents the hello request
type HelloRequest struct {
	Name string `json:"name"`
}

// HelloResponse represents the hello response
type HelloResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Err     error  `json:"err,omitempty"`
}

// Set collects all of the endpoints that compose a hello service.
type Set struct {
	HelloEndpoint endpoint.Endpoint
}

// NewSet returns a Set that wraps the provided service.
func NewSet(
	helloUC domain.HelloUseCase,
) Set {
	return Set{
		HelloEndpoint: makeHelloEndpoint(helloUC),
	}
}

func makeHelloEndpoint(uc domain.HelloUseCase) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(HelloRequest)
		response, err := uc.Execute(req.Name)
		if err != nil {
			return HelloResponse{
				Message: "Error processing request",
				Success: false,
				Err:     err,
			}, nil
		}
		return HelloResponse{
			Message: response.Message,
			Success: response.Success,
		}, nil
	}
}

// Failer is an interface that should be implemented by response types.
type Failer interface {
	Failed() error
}
