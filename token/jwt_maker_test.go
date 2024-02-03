package token

import (
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker_Token(t *testing.T) {
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
				maker, err := NewJWTMaker(secretKey)
				require.Nil(t, err)

				token, err := maker.CreateToken(args.Username, duration)
				require.Nil(t, err)
				require.NotEmpty(t, token)

				payload, err := maker.VerifyToken(token)
				require.Nil(t, err)

				require.NotEmpty(t, payload)
				require.Equal(t, args.Username, payload.Username)
				require.WithinDuration(t, args.IssuedAt, payload.IssuedAt, time.Second)
				require.WithinDuration(t, args.ExpiredAt, payload.ExpiredAt, time.Second)
			},
		},
		{
			name:      "Case 2: Invalid key size",
			secretKey: utils.RandomString(10),
			args:      &Payload{},
			testFunc: func(t *testing.T, secretKey string, args *Payload) {
				_, err := NewJWTMaker(secretKey)
				require.Error(t, err)
				require.Equal(t, InvalidKeySize, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.secretKey, tt.args)
		})
	}
}
