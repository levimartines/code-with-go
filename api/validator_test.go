package api

import (
	"code-with-go/util"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApi_Validator(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("currency", validateCurrency)
	require.NoError(t, err)

	err = validate.Var(util.USD, "currency")
	require.NoError(t, err)

	err = validate.Var("BRL", "currency")
	require.Error(t, err)

	err = validate.Var(1, "currency")
	require.Error(t, err)
}
