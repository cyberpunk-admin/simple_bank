package token

import "time"

// Maker is an interface for mamaging tokens
type Maker interface {
	// CreateToken create a new token for specified username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken check if a token is valid or not
	VerifyToken(token string) (*Payload, error)
}