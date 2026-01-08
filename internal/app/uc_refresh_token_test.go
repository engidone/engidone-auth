package app

import (
	"context"
	"testing"

	"engidoneauth/internal/jwt"
	pb "engidoneauth/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockJWTUseCase struct {
	RefreshTokenFunc func(token, refreshToken string) (*jwt.TokenInfo, error)
}

func (m *MockJWTUseCase) RefreshToken(token, refreshToken string) (*jwt.TokenInfo, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(token, refreshToken)
	}
	return &jwt.TokenInfo{}, nil
}

func TestRefreshToken_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.RefreshTokenRequest
		mockResponse   *jwt.TokenInfo
		mockError      error
		expectedCode   codes.Code
		expectedError  bool
		expectedToken  string
		expectedRefTok string
	}{
		{
			name: "successful token refresh",
			request: &pb.RefreshTokenRequest{
				Token:        "valid-token",
				RefreshToken: "valid-refresh-token",
			},
			mockResponse: &jwt.TokenInfo{
				Token:        "new-token",
				RefreshToken: "new-refresh-token",
			},
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedError: false,
			expectedToken: "new-token",
			expectedRefTok: "new-refresh-token",
		},
		{
			name: "invalid token error",
			request: &pb.RefreshTokenRequest{
				Token:        "invalid-token",
				RefreshToken: "valid-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.InvalidToken,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name: "invalid refresh token error",
			request: &pb.RefreshTokenRequest{
				Token:        "valid-token",
				RefreshToken: "invalid-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.InvalidRefreshToken,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name: "token expired error",
			request: &pb.RefreshTokenRequest{
				Token:        "expired-token",
				RefreshToken: "valid-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.TokenExpired,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name: "token generation failed error",
			request: &pb.RefreshTokenRequest{
				Token:        "valid-token",
				RefreshToken: "valid-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.TokenGenerationFailed,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
		{
			name: "refresh token generation failed error",
			request: &pb.RefreshTokenRequest{
				Token:        "valid-token",
				RefreshToken: "valid-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.GeneratingRefreshToken,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
		{
			name: "generic error fallback",
			request: &pb.RefreshTokenRequest{
				Token:        "valid-token",
				RefreshToken: "valid-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.ParsingToken,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
		{
			name: "empty tokens",
			request: &pb.RefreshTokenRequest{
				Token:        "",
				RefreshToken: "",
			},
			mockResponse:  nil,
			mockError:     jwt.InvalidToken,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name: "refresh token not found",
			request: &pb.RefreshTokenRequest{
				Token:        "valid-token",
				RefreshToken: "non-existent-refresh-token",
			},
			mockResponse:  nil,
			mockError:     jwt.RefreshTokenNotFound,
			expectedCode:  codes.NotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := &MockJWTUseCase{
				RefreshTokenFunc: func(token, refreshToken string) (*jwt.TokenInfo, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			appUC := &AppUseCase{
				jwtUC: mockJWT,
			}

			ctx := context.Background()

			resp, err := appUC.RefreshToken(ctx, tt.request)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("expected gRPC status error, got: %v", err)
					return
				}

				if st.Code() != tt.expectedCode {
					t.Errorf("expected code %v, got %v", tt.expectedCode, st.Code())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if resp == nil {
					t.Errorf("expected response but got nil")
					return
				}

				if resp.Token != tt.expectedToken {
					t.Errorf("expected token %q, got %q", tt.expectedToken, resp.Token)
				}

				if resp.RefreshToken != tt.expectedRefTok {
					t.Errorf("expected refresh token %q, got %q", tt.expectedRefTok, resp.RefreshToken)
				}
			}
		})
	}
}
