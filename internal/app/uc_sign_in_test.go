package app

import (
	"context"
	"errors"
	"testing"

	"engidoneauth/internal/credentials"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
	"engidoneauth/internal/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockSignInUseCase struct {
	SingInFunc func(credentials signin.Credentials) (*signin.Result, error)
}

func (m *MockSignInUseCase) SingIn(credentials signin.Credentials) (*signin.Result, error) {
	if m.SingInFunc != nil {
		return m.SingInFunc(credentials)
	}
	return &signin.Result{}, nil
}

type MockableError struct {
	Err error
}

func (e *MockableError) Error() string {
	return e.Err.Error()
}

func (e *MockableError) Unwrap() error {
	return e.Err
}

var (
	MockInvalidCredentials     = &MockableError{Err: credentials.InvalidCredentials}
	MockUserNotFound           = &MockableError{Err: users.UserNotFound}
	MockMissingUsername        = &MockableError{Err: signin.MissingUsername}
	MockMissingPassword        = &MockableError{Err: signin.MissingPassword}
	MockUsernameTooShort       = &MockableError{Err: signin.UsernameTooShort}
	MockPasswordTooShort       = &MockableError{Err: signin.PasswordTooShort}
	MockCredentialsInternal    = &MockableError{Err: credentials.InternalError}
	MockCredentialsUserNotFound = &MockableError{Err: credentials.UserNotFound}
)

func TestSignIn_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.SignInRequest
		mockResponse   *signin.Result
		mockError      error
		expectedCode   codes.Code
		expectedError  bool
		expectedToken  string
		expectedRefTok string
	}{
		{
			name: "successful sign in",
			request: &pb.SignInRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockResponse: &signin.Result{
				Token:        "generated-token",
				RefreshToken: "generated-refresh-token",
			},
			mockError:      nil,
			expectedCode:   codes.OK,
			expectedError:  false,
			expectedToken:  "generated-token",
			expectedRefTok: "generated-refresh-token",
		},
		{
			name: "invalid credentials error",
			request: &pb.SignInRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockResponse:  nil,
			mockError:     MockInvalidCredentials,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name: "user not found error",
			request: &pb.SignInRequest{
				Username: "nonexistent",
				Password: "password",
			},
			mockResponse:  nil,
			mockError:     MockUserNotFound,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name: "missing username error",
			request: &pb.SignInRequest{
				Username: "",
				Password: "password",
			},
			mockResponse:  nil,
			mockError:     MockMissingUsername,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name: "missing password error",
			request: &pb.SignInRequest{
				Username: "testuser",
				Password: "",
			},
			mockResponse:  nil,
			mockError:     MockMissingPassword,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name: "username too short error",
			request: &pb.SignInRequest{
				Username: "ab",
				Password: "password",
			},
			mockResponse:  nil,
			mockError:     MockUsernameTooShort,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name: "password too short error",
			request: &pb.SignInRequest{
				Username: "testuser",
				Password: "abc",
			},
			mockResponse:  nil,
			mockError:     MockPasswordTooShort,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name: "generic error fallback",
			request: &pb.SignInRequest{
				Username: "testuser",
				Password: "password",
			},
			mockResponse:  nil,
			mockError:     MockCredentialsInternal,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
		{
			name: "both username and password empty",
			request: &pb.SignInRequest{
				Username: "",
				Password: "",
			},
			mockResponse:  nil,
			mockError:     MockMissingUsername,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name: "valid credentials with user not found",
			request: &pb.SignInRequest{
				Username: "testuser",
				Password: "password",
			},
			mockResponse:  nil,
			mockError:     MockCredentialsUserNotFound,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSignIn := &MockSignInUseCase{
				SingInFunc: func(credentials signin.Credentials) (*signin.Result, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			appUC := &AppUseCase{
				signInUC: mockSignIn,
			}

			ctx := context.Background()

			resp, err := appUC.SignIn(ctx, tt.request)

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
					t.Logf("Expected code: %v, got: %v", tt.expectedCode, st.Code())
					t.Logf("Error message: %v", st.Message())
					t.Logf("Mock error type: %T", tt.mockError)
					if unwrapErr := errors.Unwrap(tt.mockError); unwrapErr != nil {
						t.Logf("Unwrapped mock error type: %T", unwrapErr)
					}
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
