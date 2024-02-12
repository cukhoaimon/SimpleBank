package token

import (
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker_Token(t *testing.T) {
	duration := time.Hour
	payloadId, err := uuid.NewUUID()
	require.Nil(t, err)

	tests := []struct {
		name      string
		secretKey string
		args      *Payload
		testFunc  func(t *testing.T, secretKey string, args *Payload)
	}{
		{
			name:      "Case 1: Happy case",
			secretKey: utils.RandomString(32),
			args: &Payload{
				Id:        payloadId,
				Username:  utils.RandomOwner(),
				IssuedAt:  time.Now(),
				ExpiredAt: time.Now().Add(duration),
			},
			testFunc: func(t *testing.T, secretKey string, args *Payload) {
				maker, err := NewPasetoMaker(secretKey)
				require.Nil(t, err)

				token, returnedPayload, err := maker.CreateToken(args.Username, duration)
				require.Nil(t, err)
				require.NotEmpty(t, token)
				verifyPayload(t, args, returnedPayload)

				payload, err := maker.VerifyToken(token)
				require.Nil(t, err)
				verifyPayload(t, args, payload)
			},
		},
		{
			name:      "Case 2: Invalid key size",
			secretKey: utils.RandomString(10),
			args:      &Payload{},
			testFunc: func(t *testing.T, secretKey string, args *Payload) {
				_, err := NewPasetoMaker(secretKey)
				require.Error(t, err)
				require.Equal(t, PasetoInvalidKeySize, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.secretKey, tt.args)
		})
	}
}

func verifyPayload(t *testing.T, want *Payload, have *Payload) {
	require.NotEmpty(t, have)
	require.Equal(t, want.Username, have.Username)
	require.WithinDuration(t, want.IssuedAt, have.IssuedAt, time.Second)
	require.WithinDuration(t, want.ExpiredAt, have.ExpiredAt, time.Second)
}
