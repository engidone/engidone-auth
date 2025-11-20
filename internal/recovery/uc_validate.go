package recovery

import (
	"engidoneauth/log"
)
func (uc *UseCase) ValidateToken(code string) (*ValidateResponse, error) {
	recoveryCode, err := uc.repo.findRecoveryCode(code)
	if err != nil {
		log.Error("Error finding recovery code", log.Err(err), log.String("code", code))
		return nil, err
	}

	if recoveryCode != code {
		log.Warn("Invalid recovery code", log.String("code", code))
		return &ValidateResponse{Success: false, Message: "Invalid code"}, nil
	}

	return &ValidateResponse{Success: true, Message: "Token is valid"}, nil
}