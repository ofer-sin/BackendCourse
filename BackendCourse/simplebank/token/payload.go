package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredTocken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

// Payload represents the data stored in a token.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a unique ID, the specified username,
// and sets the creation and expiration times based on the provided duration.
// Returns the created Payload or an error if the token ID could not be generated.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredTocken
	}

	return nil
}
