package app

import "engidoneauth/internal/signin"

type SignInResponse struct {
	signin.Result
	Success bool
	Message string
}
