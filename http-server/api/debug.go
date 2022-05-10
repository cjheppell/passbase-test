package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cjheppell/passbase-test/http-server/user"
)

type debugRepo interface {
	GetAllUsers() []user.User
}

type debugHandler struct {
	repo debugRepo
}

func NewDebugHandler(r debugRepo) debugHandler {
	return debugHandler{
		repo: r,
	}
}

func (h debugHandler) PrintAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.repo.GetAllUsers()
	usersJson, err := json.Marshal(users)
	if err != nil {
		fmt.Fprintf(w, "error marshalling users: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", string(usersJson))
}
