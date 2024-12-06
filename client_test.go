package unifiaccessclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/miquelruiz/go-unifi-access-api/schema"
)

var (
	token          = "some-auth-token"
	firstName      = "FirstName"
	lastName       = "LastName"
	employeeNumber = "123"
	email          = "useremail@example.com"
	onboardTime    = int(time.Now().Unix())
	expectedId     = "9876"
)

func TestFailedNew(t *testing.T) {
	_, err := New(":some@nasty/malformed\\stuff", token)
	if err == nil {
		t.Errorf("Error was expected, constructor didn't fail")
	}
}

func TestCreateUser(t *testing.T) {
	expectedPath := "/api/v1/developer/users"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPath {
			t.Errorf("Unexpected path called: %s. Expected: %s", r.URL.Path, expectedPath)
		}

		w.WriteHeader(http.StatusOK)
		rawResponse, err := json.Marshal(schema.Response[schema.UserResponse]{
			Code: schema.ResponseSuccess,
			Msg:  "success",
			Data: schema.UserResponse{
				Id: expectedId,
			},
		})
		if err != nil {
			t.Errorf("Error marshaling mock response: %v", err)
		}
		w.Write(rawResponse)
	}))

	client, err := New(server.URL, token)
	if err != nil {
		t.Errorf("Unexpected error building client: %v", err)
	}

	resp, err := client.CreateUser(schema.UserRequest{
		FirstName:      firstName,
		LastName:       lastName,
		UserEmail:      &email,
		EmployeeNumber: &employeeNumber,
		OnboardTime:    &onboardTime,
	})
	if err != nil {
		t.Errorf("Unexpected failure calling CreateUser: %v", err)
	}

	if resp.Id != expectedId {
		t.Errorf("Unexpected ID received: %s. Expected: %s", resp.Id, expectedId)
	}
}

func TestGetUser(t *testing.T) {
	expectedPath := fmt.Sprintf("/api/v1/developer/users/%s", expectedId)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPath {
			t.Errorf("Unexpected path called: %s. Expected: %s", r.URL.Path, expectedPath)
		}

		w.WriteHeader(http.StatusOK)
		rawResponse, err := json.Marshal(schema.Response[schema.UserResponse]{
			Code: schema.ResponseSuccess,
			Msg:  "success",
			Data: schema.UserResponse{
				Id: expectedId,
			},
		})
		if err != nil {
			t.Errorf("Error marshaling mock response: %v", err)
		}
		w.Write(rawResponse)

	}))

	client, err := New(server.URL, token)
	if err != nil {
		t.Errorf("Unexpected error creating client: %v", err)
	}

	user, err := client.GetUser(expectedId)
	if err != nil {
		t.Errorf("Unexpected error calling GetUser: %v", err)
	}

	if user.Id != expectedId {
		t.Errorf("Unexpected Id from GetUser: %s. Expected: %s", user.Id, expectedId)
	}
}
