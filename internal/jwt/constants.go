package jwt

const (
	ErrInvalidToken               = "INVALID_TOKEN"
	ErrLoadingPublicKey           = "ERROR_LOADING_PUBLIC_KEY"
	ErrLoadingPrivateKey          = "ERROR_LOADING_PUBLIC_KEY"
	ErrGeneratingRefreshToken     = "ERROR_GENERATING_REFRESH_TOKEN"
	ErrParsingToken               = "ERROR_PARSING_TOKEN"
	ErrInsertOrUpdateRefreshToken = "ERROR_INSERT_OR_UPDATE_TOKEN"
	ErrInvalidRefreshToken        = "INVALID_REFRESH_TOKEN"
)
