// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr implementation of the UserService
// ----------------------------------------------------------------------------

package impl

import (
	"encoding/json"
	"os"

	"github.com/benc-uk/dapr-store/cmd/users/spec"
	"github.com/benc-uk/dapr-store/pkg/dapr"
	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// UserService is a Dapr based implementation of UserService interface
type UserService struct {
	*dapr.Helper
	storeName string
}

// NewService creates a new UserService
func NewService(serviceName string) *UserService {
	// Set up Dapr & checks for Dapr sidecar port, abort
	helper := dapr.NewHelper(serviceName)
	if helper == nil {
		os.Exit(1)
	}
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	return &UserService{
		helper,
		storeName,
	}
}

// AddUser registers a new user and stores in Dapr state
func (s *UserService) AddUser(user spec.User) error {
	// Check is user already registered
	data, prob := s.GetState(s.storeName, user.Username)
	if prob != nil {
		return prob
	}

	// If we get any data, that means we found a user, that's an error in our case
	if len(data) > 0 {
		prob := problem.New("err://user-exists", user.Username+" already registered", 400, user.Username+" already registered", s.ServiceName)
		return prob
	}

	// Call Dapr helper to save state
	prob = s.SaveState(s.storeName, user.Username, user)
	if prob != nil {
		return prob
	}

	return nil
}

// GetUser fetches a user from Dapr state
func (s *UserService) GetUser(username string) (*spec.User, error) {
	data, prob := s.GetState(s.storeName, username)
	if prob != nil {
		return nil, prob
	}

	if len(data) <= 0 {
		prob := problem.New("err://not-found", "No data returned", 404, "Username: '"+username+"' not found", s.ServiceName)
		return nil, prob
	}

	user := &spec.User{}
	err := json.Unmarshal(data, user)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed user JSON", 500, "JSON could not be decoded", s.ServiceName)
		return nil, prob
	}

	return user, nil
}
