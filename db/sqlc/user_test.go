package db

import (
	"context"
	"github.com/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) User {
	password, err := util.HashPassword(util.RandomString(6))

	require.NoError(t, err)
	require.NotEmpty(t, password)

	arg := CreateUserParams{
		UserName:     util.RandomOwner(),
		HashPassword: password,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.UserName, arg.UserName)
	require.Equal(t, user.HashPassword, arg.HashPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)

	require.True(t, user.PasswordChangeAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	_ = CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)

}
