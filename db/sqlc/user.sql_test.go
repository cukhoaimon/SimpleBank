package db

import (
	"context"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: "secret",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	want := createRandomUser(t)

	have, err := testQuery.GetUser(context.Background(), want.Username)

	require.Nil(t, err)
	require.Equal(t, want, have)
}
