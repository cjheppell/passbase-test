package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cjheppell/passbase/test-app/http-server/user"
)

type userRepo interface {
	CreateIfNotExist(userId user.UserId) (user.User, error)
}

type BasicAuth struct {
	handler  http.Handler
	userRepo userRepo
}

var userContextKey = "userContextKey"

func (a BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, _, ok := r.BasicAuth()
	if !ok {
		// prompt to provide basic auth credentials
		w.Header().Set("WWW-Authenticate", `Basic realm="/", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// we'll just add a new user every time we see them
	// this is clearly a bit pointless, but it gives us the concept of a "user"
	userId := user.UserId(username)

	u, err := a.userRepo.CreateIfNotExist(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), userContextKey, u))

	a.handler.ServeHTTP(w, r)
}

func NewBasicAuth(handlerToWrap http.Handler, repo userRepo) *BasicAuth {
	return &BasicAuth{handlerToWrap, repo}
}

func GetUserFromContext(ctx context.Context) (*user.User, error) {
	u := ctx.Value(userContextKey)
	if u == nil {
		return nil, fmt.Errorf("failed to find user id in context")
	}
	foundUser := u.(user.User)
	return &foundUser, nil
}
