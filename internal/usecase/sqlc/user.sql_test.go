package usecase

import (
	"context"
	"database/sql"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(10))
	require.Nil(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
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

func TestQueries_UpdateUser(t *testing.T) {
	oldUser := createRandomUser(t)

	newHashedPassword, err := utils.HashPassword("secret")
	require.Nil(t, err)

	tests := []struct {
		name       string
		params     UpdateUserParams
		assertFunc func(t *testing.T, updated User, params UpdateUserParams)
	}{
		{
			name: "Update email",
			params: UpdateUserParams{
				Email:    newSqlString(utils.RandomEmail()),
				Username: oldUser.Username,
			},
			assertFunc: func(t *testing.T, updated User, params UpdateUserParams) {
				require.Equal(t, params.Email.String, updated.Email)
				require.Equal(t, oldUser.FullName, updated.FullName)
				require.Equal(t, oldUser.Username, updated.Username)
			},
		},
		{
			name: "Update password",
			params: UpdateUserParams{
				// new password is secret
				HashedPassword: newSqlString(newHashedPassword),
				Username:       oldUser.Username,
			},
			assertFunc: func(t *testing.T, updated User, params UpdateUserParams) {
				require.Nil(t, utils.CheckPassword("secret", updated.HashedPassword))
				require.Equal(t, oldUser.Email, updated.Email)
				require.Equal(t, oldUser.FullName, updated.FullName)
				require.Equal(t, oldUser.Username, updated.Username)
			},
		},
		{
			name: "Update full name",
			params: UpdateUserParams{
				FullName: newSqlString(utils.RandomOwner()),
				Username: oldUser.Username,
			},
			assertFunc: func(t *testing.T, updated User, params UpdateUserParams) {
				require.Equal(t, params.FullName.String, updated.FullName)
				require.Equal(t, oldUser.Email, updated.Email)
				require.Equal(t, oldUser.Username, updated.Username)
			},
		},
	}

	for _, tt := range tests {
		go t.Run(tt.name, func(t *testing.T) {
			updatedUser, err := testQuery.UpdateUser(context.Background(), tt.params)
			require.Nil(t, err)
			require.NotEmpty(t, updatedUser)

			tt.assertFunc(t, updatedUser, tt.params)
		})
	}
}

func newSqlString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}
