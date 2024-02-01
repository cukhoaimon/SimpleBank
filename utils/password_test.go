package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)
	hashedPassword, err := HashPassword(password)

	require.Nil(t, err)

	type args struct {
		password       string
		hashedPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Same password",
			args: args{
				password:       password,
				hashedPassword: hashedPassword,
			},
			wantErr: false,
		},
		{
			name: "Different password",
			args: args{
				password:       RandomString(10),
				hashedPassword: hashedPassword,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPassword(tt.args.password, tt.args.hashedPassword); (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
