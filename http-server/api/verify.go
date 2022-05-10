package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cjheppell/passbase/test-app/http-server/auth"
	"github.com/cjheppell/passbase/test-app/http-server/user"
)

type userIdVerificationRepository interface {
	AssociatePassbaseKey(userId user.UserId, passbaseKey string) error
}

type verifyHandler struct {
	repo userIdVerificationRepository
}

func NewVerifyHandler(repo userIdVerificationRepository) verifyHandler {
	return verifyHandler{
		repo: repo,
	}
}

const verifyHtml = `
<html>
  <head>
    <title>Verify with Passbase</title>
		<script
			type="text/javascript"
			src="https://unpkg.com/@passbase/button"
		></script>
  </head>
  <body>
    <h1>Verify</h1>
    <div>
      <p>Verify your identity by clicking the button below</p>
      <!-- Place the code below where you want your button to appear -->
      <div id="passbase-button"></div>
      <script type="text/javascript">
        Passbase.renderButton(
          document.getElementById("passbase-button"),
          "eav0UeDt7ZF5IDaPlPHhQBP89ieM9sgE9vMxkwAL0HYf94bSdtqxCKuQPFcPtJnd",
          {
            onStart: () => {},
            onError: (errorCode) => {},
            onFinish: (identityAccessKey) => {fetch('/verify', {method: 'POST', headers: { 'Content-Type': 'application/json'}, body: JSON.stringify({key: identityAccessKey})} ).catch(err => {console.error(err)})},
          }
        );
      </script>
  </body>
</html>
`

func (v verifyHandler) RenderVerify(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, verifyHtml)
}

func (v verifyHandler) BeginVerify(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Key string `json:"key"`
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading payload")
		return
	}

	p := payload{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := v.repo.AssociatePassbaseKey(user.Id, p.Key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to associate user with passbase verification key")
		return
	}

	fmt.Printf("successfully associated passbase key %q with user %q\n", p.Key, user.Id)

	w.WriteHeader(http.StatusOK)
}
