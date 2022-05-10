package api

import (
	"fmt"
	"net/http"
)

type VerifiedOnlyHandler struct {
}

func (h VerifiedOnlyHandler) SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello from the verified only handler\nIf you're reading this, then Passbase verified you successfully!")
}
