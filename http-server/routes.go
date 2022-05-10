package main

import (
	"net/http"

	"github.com/cjheppell/passbase/test-app/http-server/api"
	"github.com/cjheppell/passbase/test-app/http-server/auth"
	"github.com/cjheppell/passbase/test-app/http-server/user"
	"github.com/cjheppell/passbase/test-app/http-server/webhook"
)

func register(mux *http.ServeMux) {
	userRepo := user.NewUserRepository()

	pbWebhook := webhook.NewPassbaseWebhookHandler(&userRepo)
	mux.HandleFunc("/passbase/event", pbWebhook.ReceiveWebhookEvent)

	verifyHandler := api.NewVerifyHandler(&userRepo)
	verifyMethodRouter := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			verifyHandler.BeginVerify(rw, r)
		}
		if r.Method == http.MethodGet {
			verifyHandler.RenderVerify(rw, r)
		}
	})
	mux.Handle("/verify", auth.NewBasicAuth(verifyMethodRouter, &userRepo))

	verifOnlyHandler := api.VerifiedOnlyHandler{}
	mux.Handle("/verified", auth.NewBasicAuth(auth.NewIdVerifiedAuthorization((http.HandlerFunc(verifOnlyHandler.SayHello))), &userRepo))

	debugHandler := api.NewDebugHandler(&userRepo)
	mux.HandleFunc("/debug", debugHandler.PrintAllUsers)
}
