package handler

import (
	"encoding/json"
	"net/http"

	server "github.com/ykpythemind/apispec-handler-linter/apispec"
)

// Handlers implements ServerInterface
type Handlers struct {
}

func NewHandlers() server.ServerInterface {
	return &Handlers{}
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	var _ server.GetUsersRequest // 不要だが使わなければいけない...

	// do something...

	res := server.GetUsersResponse{
		{Email: "test@example.com"},
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req server.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// do something...

	res := server.CreateUserResponse{
		User: &server.User{},
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
