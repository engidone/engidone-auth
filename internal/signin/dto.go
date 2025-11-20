package signin

import "time"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type SigInResponse struct {
	Token    string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}