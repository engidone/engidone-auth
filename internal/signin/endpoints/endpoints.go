package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"engidone-auth/internal/signin/domain"
)

// SigninRequest represents the signin request
type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SigninResponse represents the signin response
type SigninResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Token    string `json:"token,omitempty"`
	ExpiresAt int64 `json:"expires_at,omitempty"`
	Err      error  `json:"err,omitempty"`
}

// ValidateTokenRequest represents the validate token request
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// ValidateTokenResponse represents the validate token response
type ValidateTokenResponse struct {
	Valid    bool   `json:"valid"`
	Message  string `json:"message"`
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Err      error  `json:"err,omitempty"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// GetUserRequest represents the get user request
type GetUserRequest struct {
	UserID string `json:"user_id"`
}

// GetUserResponse represents the get user response
type GetUserResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
	Err      error  `json:"err,omitempty"`
}

var logger log.Logger

// Set collects all of the endpoints that compose an auth service.
type Set struct {
	SigninEndpoint       endpoint.Endpoint
	ValidateTokenEndpoint endpoint.Endpoint
	RefreshTokenEndpoint  endpoint.Endpoint
	GetUserEndpoint       endpoint.Endpoint
}

// NewSet returns a Set that wraps the provided server.
func NewSet(
	signinUC domain.SigninUseCase,
	validateTokenUC domain.ValidateTokenUseCase,
	refreshTokenUC domain.RefreshTokenUseCase,
	getUserUC domain.GetUserUseCase,
	log log.Logger,
) Set {
	logger = log

	return Set{
		SigninEndpoint:       makeSigninEndpoint(signinUC),
		ValidateTokenEndpoint: makeValidateTokenEndpoint(validateTokenUC),
		RefreshTokenEndpoint:  makeRefreshTokenEndpoint(refreshTokenUC),
		GetUserEndpoint:       makeGetUserEndpoint(getUserUC),
	}
}

func makeSigninEndpoint(uc domain.SigninUseCase) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SigninRequest)
		credentials := domain.Credentials{
			Username: req.Username,
			Password: req.Password,
		}
		authResponse, err := uc.Execute(credentials)
		if err != nil {
			return SigninResponse{
				Success: false,
				Message: "Authentication failed",
				Err:     err,
			}, nil
		}
		return SigninResponse{
			Success:   true,
			Message:   "Authentication successful",
			UserID:    authResponse.UserID,
			Username:  authResponse.Username,
			Email:     authResponse.Email,
			Token:     authResponse.Token,
			ExpiresAt: authResponse.ExpiresAt.Unix(),
		}, nil
	}
}

func makeValidateTokenEndpoint(uc domain.ValidateTokenUseCase) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateTokenRequest)
		user, err := uc.Execute(req.Token)
		if err != nil {
			return ValidateTokenResponse{
				Valid:   false,
				Message: "Invalid token",
				Err:     err,
			}, nil
		}
		return ValidateTokenResponse{
			Valid:    true,
			Message:  "Token valid",
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
		}, nil
	}
}

func makeRefreshTokenEndpoint(uc domain.RefreshTokenUseCase) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RefreshTokenRequest)
		authResponse, err := uc.Execute(req.UserID, req.Token)
		if err != nil {
			return SigninResponse{
				Success: false,
				Message: "Token refresh failed",
				Err:     err,
			}, nil
		}
		return SigninResponse{
			Success:   true,
			Message:   "Token refreshed successfully",
			UserID:    authResponse.UserID,
			Username:  authResponse.Username,
			Email:     authResponse.Email,
			Token:     authResponse.Token,
			ExpiresAt: authResponse.ExpiresAt.Unix(),
		}, nil
	}
}

func makeGetUserEndpoint(uc domain.GetUserUseCase) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := uc.Execute(req.UserID)
		if err != nil {
			return GetUserResponse{
				Success: false,
				Message: "User not found",
				Err:     err,
			}, nil
		}
		return GetUserResponse{
			Success:   true,
			Message:   "User found",
			UserID:    user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		}, nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response types may implement Failed() error method to indicate
// if their response should be considered as failure.
type Failer interface {
	Failed() error
}