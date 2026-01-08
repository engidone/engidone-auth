package signin

type Credentials struct {
	Username string
	Password string
}

type Result struct {
	Token        string
	RefreshToken string
}
