package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/simplebank/db/mock"
	db "github.com/simplebank/db/sqlc"
	util2 "github.com/simplebank/util"
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

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util2.CheckPassword(e.password, arg.HashPassword)
	if err != nil {
		return false
	}
	e.arg.HashPassword = arg.HashPassword
	return reflect.DeepEqual(e.arg, arg)
}
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %s", e.arg, e.password)
}

func eqCreateUserParams(params db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{params, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := RandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.UserName,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					UserName:     user.UserName,
					HashPassword: password,
					FullName:     user.FullName,
					Email:        user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"user_name": user.UserName,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
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
		{
			name: "DuplicateUserName",
			body: gin.H{
				"user_name": user.UserName,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUserName",
			body: gin.H{
				"user_name": "invai$ld--$$#1",
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"user_name": user.UserName,
				"full_name": user.FullName,
				"email":     "@gmail.xxx.yyyzz",
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToShortPassword",
			body: gin.H{
				"user_name": "invai$ld--$$#1",
				"full_name": user.UserName,
				"email":     user.Email,
				"password":  "mm",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start the test server and request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			// marshal data to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

//func TestLoginUserAPI(t *testing.T) {
//	user, password := RandomUser(t)
//
//	testCases := []struct {
//		name          string
//		body          gin.H
//		buildStubs    func(store *mockdb.MockStore)
//		checkResponse func(recoder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "OK",
//			body: gin.H{
//				"username": user.UserName,
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Eq(user.UserName)).
//					Times(1).
//					Return(user, nil)
//				store.EXPECT().
//					CreateSession(gomock.Any(), gomock.Any()).
//					Times(1)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//			},
//		},
//		{
//			name: "UserNotFound",
//			body: gin.H{
//				"username": "NotFound",
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, sql.ErrNoRows)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusNotFound, recorder.Code)
//			},
//		},
//		{
//			name: "IncorrectPassword",
//			body: gin.H{
//				"username": user.UserName,
//				"password": "incorrect",
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Eq(user.UserName)).
//					Times(1).
//					Return(user, nil)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//			},
//		},
//		{
//			name: "InternalError",
//			body: gin.H{
//				"username": user.UserName,
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, sql.ErrConnDone)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusInternalServerError, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidUsername",
//			body: gin.H{
//				"username":  "invalid-user#1",
//				"password":  password,
//				"full_name": user.FullName,
//				"email":     user.Email,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//			server := NewTestServer(t, store)
//			recorder := httptest.NewRecorder()
//
//			// Marshal body data to JSON
//			data, err := json.Marshal(tc.body)
//			require.NoError(t, err)
//
//			url := "/users/login"
//			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
//			require.NoError(t, err)
//
//			server.router.ServeHTTP(recorder, request)
//			tc.checkResponse(recorder)
//		})
//	}
//}

func RandomUser(t *testing.T) (user db.User, password string) {
	password = util2.RandomString(6)
	hashPassword, err := util2.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		UserName:     util2.RandomOwner(),
		FullName:     util2.RandomOwner(),
		HashPassword: hashPassword,
		Email:        util2.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var getUser db.User
	err = json.Unmarshal(data, &getUser)

	require.NoError(t, err)
	require.Equal(t, user.UserName, getUser.UserName)
	require.Equal(t, user.FullName, getUser.FullName)
	require.Equal(t, user.Email, getUser.Email)
	require.Empty(t, getUser.HashPassword)
}
