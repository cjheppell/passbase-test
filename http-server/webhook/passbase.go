package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cjheppell/passbase/test-app/http-server/user"
)

type passbaseEvent string

const (
	verificationComplete  passbaseEvent = "VERIFICATION_COMPLETED"
	verificationReviewed  passbaseEvent = "VERIFICATION_REVIEWED"
	datapointUpdated      passbaseEvent = "DATAPOINT_UPDATED"
	watchlistMonitoring   passbaseEvent = "WATCHLIST_MONITORING"
	identityAuthenticated passbaseEvent = "IDENTITY_AUTHENTICATED"
)

type passbaseWebhookEventPayload struct {
	EventType passbaseEvent `json:"event"`
	Key       string        `json:"key"`
	Created   int           `json:"created"`
	Updated   int           `json:"updated"`
}

type verificationCompletedPayload struct {
	passbaseWebhookEventPayload
	Status    string `json:"status"`
	Processed int    `json:"processed"`
}

type verificationReviewedPayload struct {
	verificationCompletedPayload
}

type datapointUpdatedPayload struct {
	passbaseWebhookEventPayload
	ResourceKey string `json:"resource_key"`
	Types       string `json:"type"`
	Value       bool   `json:"value"`
}

type watchlistMonitoringPayload struct {
	passbaseWebhookEventPayload
	Types []string `json:"types"`
	Clean bool     `json:"clean"`
}

type identityAuthenticatedPayload struct {
	passbaseWebhookEventPayload
	Status    string `json:"status"`
	Processed int    `json:"processed"`
}

type idVerificationRepository interface {
	GetUserFromPassbaseKey(passbaseKey string) (*user.User, error)
	RegisterUserVerified(userId user.UserId) error
}

type PassbaseWebhookHandler struct {
	repo idVerificationRepository
}

func NewPassbaseWebhookHandler(repo idVerificationRepository) PassbaseWebhookHandler {
	return PassbaseWebhookHandler{
		repo: repo,
	}
}

func (p *PassbaseWebhookHandler) ReceiveWebhookEvent(w http.ResponseWriter, r *http.Request) {
	bodyContents, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = p.handleEventPayload(bodyContents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error printing webhook event payload: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *PassbaseWebhookHandler) handleEventPayload(bodyContents []byte) error {
	pbEvent := passbaseWebhookEventPayload{}
	if err := json.Unmarshal(bodyContents, &pbEvent); err != nil {
		return err
	}

	switch pbEvent.EventType {
	case verificationComplete:
		payload := verificationCompletedPayload{}
		if err := json.Unmarshal(bodyContents, &payload); err != nil {
			return err
		}

		user, err := p.repo.GetUserFromPassbaseKey(payload.Key)
		if err != nil {
			return fmt.Errorf("failed to find user from passbase key: %w", err)
		}

		if payload.Status == "approved" {
			err := p.repo.RegisterUserVerified(user.Id)
			if err != nil {
				return err
			}
		}

	case verificationReviewed:
		payload := verificationReviewedPayload{}
		if err := json.Unmarshal(bodyContents, &payload); err != nil {
			return err
		}
		fmt.Printf("received verif reviewed event: %+v", payload)
	case datapointUpdated:
		payload := datapointUpdatedPayload{}
		if err := json.Unmarshal(bodyContents, &payload); err != nil {
			return err
		}
		fmt.Printf("received data point updated event: %+v", payload)
	case watchlistMonitoring:
		payload := watchlistMonitoringPayload{}
		if err := json.Unmarshal(bodyContents, &payload); err != nil {
			return err
		}
		fmt.Printf("received watchlist monitoring event: %+v", payload)
	case identityAuthenticated:
		payload := identityAuthenticatedPayload{}
		if err := json.Unmarshal(bodyContents, &payload); err != nil {
			return err
		}
		fmt.Printf("received id auth event: %+v", payload)
	}

	return nil
}