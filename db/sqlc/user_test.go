package db

import (
	"code-with-go/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user := createRandomUser(t)
	retrievedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)

	require.Equal(t, retrievedUser.Username, user.Username)
	require.Equal(t, retrievedUser.HashedPassword, user.HashedPassword)
	require.Equal(t, retrievedUser.FullName, user.FullName)
	require.Equal(t, retrievedUser.Email, user.Email)
}

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	return user
}
