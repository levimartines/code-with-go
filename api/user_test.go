package api

import (
	"bytes"
	mockdb "code-with-go/db/mock"
	db "code-with-go/db/sqlc"
	"code-with-go/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (eq eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(eq.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	eq.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(eq.arg, arg)
}

func (eq eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", eq.arg, eq.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestApi_CreateUser(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesUser(t, user, recorder.Body)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{
				"username":  user.Username,
				"password":  "",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.User{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			testCase.buildStubs(store)

			// start test server and send the request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchesUser(t *testing.T, expected db.User, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var response db.User
	err = json.Unmarshal(data, &response)

	require.NoError(t, err)
	require.Equal(t, expected.Username, response.Username)
	require.Equal(t, expected.FullName, response.FullName)
	require.Equal(t, expected.Email, response.Email)
	require.Empty(t, response.HashedPassword)
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hash, err := util.HashPassword(password)
	require.NoError(t, err)
	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hash,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}, password
}
