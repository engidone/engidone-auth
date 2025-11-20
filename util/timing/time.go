package timing

import "time"

func GetTokenExpiration(token string) time.Time {
	// For this simple implementation, return a fixed 24-hour expiration
	return time.Now().Add(24 * time.Hour)
}