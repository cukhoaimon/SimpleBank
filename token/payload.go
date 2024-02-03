package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.ExpiredAt}, nil
}

func (p Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.IssuedAt}, nil
}

func (p Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.IssuedAt}, nil
}

func (p Payload) GetIssuer() (string, error) {
	return "jwt_maker.go", nil
}

func (p Payload) GetSubject() (string, error) {
	return "user payload", nil
}

func (p Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{"aud"}, nil
}

// NewPayload create a new payload with specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Id:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}
