package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ofer-sin/Courses/BackendCourse/simplebank/db/mock"
	db "github.com/ofer-sin/Courses/BackendCourse/simplebank/db/sqlc"
	"github.com/ofer-sin/Courses/BackendCourse/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t) //required for mocking

	// This line is not needed here, as we are using the controller in the test function
	// It will be called at the end of the test function
	// to clean up the mock controller
	// and release any resources it holds.
	// This is a common pattern in Go testing with gomock
	defer ctrl.Finish()

	// Create a mock store (instead of the real store)
	store := mockdb.NewMockStore(ctrl)

	// Set up the expected call to the mock store GetAccount() method
	// with the context and the account ID
	// and specify that it should be called once
	// and return the account and nil error
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// Create a new server with the mock store
	server := NewServer(store)

	// Create a new HTTP request to the /accounts/:id endpoint
	// Create a new HTTP recorder to capture the response
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Serve the HTTP request using the test server
	// This will call the getAccount() method of the server
	// which will in turn call the GetAccount() method of the mock store
	server.router.ServeHTTP(recorder, request)

	// Check the response status code, which is recorder in the recorder
	require.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Accounts
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
