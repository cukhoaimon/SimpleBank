package token

import "time"

// Maker interface for managing token
type Maker interface {
	// CreateToken create token from username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if token is valid or not
	VerifyToken(token string) (*Payload, error)
}
