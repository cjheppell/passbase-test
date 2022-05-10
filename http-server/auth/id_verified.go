package auth

import (
	"fmt"
	"net/http"
)

type IdVerifiedAuthorization struct {
	handler http.Handler
}

func (a IdVerifiedAuthorization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !user.IdentityVerified {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "User not authorized to access this service: You have not verified your identity with Passbase\n")
		return
	}

	a.handler.ServeHTTP(w, r)
}

func NewIdVerifiedAuthorization(handlerToWrap http.Handler) *IdVerifiedAuthorization {
	return &IdVerifiedAuthorization{handlerToWrap}
}
