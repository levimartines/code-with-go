package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_Password(t *testing.T) {
	rawPassword := RandomString(6)

	hash, err := HashPassword(rawPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	newPassword, err := HashPassword(rawPassword)
	require.NoError(t, err)
	require.NotEqual(t, hash, newPassword)

	err = CheckPassword(rawPassword, hash)
	require.NoError(t, err)

	err = CheckPassword("hash", hash)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
