package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_Password(t *testing.T) {
	rawPassword := RandomString(6)

	password, err := HashPassword(rawPassword)
	require.NoError(t, err)
	require.NotEmpty(t, password)

	newPassword, err := HashPassword(rawPassword)
	require.NoError(t, err)
	require.NotEqual(t, password, newPassword)

	err = CheckPassword(password, rawPassword)
	require.NoError(t, err)

	err = CheckPassword(password, "password")
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
