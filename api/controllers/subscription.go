package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/payment"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// TODO: Check user is owner of organization
func PostSubscription(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	_ models.User,
) (response router.Serde, status int, err error) {
	status = http.StatusOK
	var sessionID string
	var url string

	p := payment.NewStripePayment()

	organizationID, _ := strconv.ParseUint(
		params.Get("organizationID").(string),
		10,
		64,
	)

	orga := models.Organization{ID: uint(organizationID)}

	if err = Repo.GetOrganization(&orga).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		goto done
	}

	sessionID, url, err = p.StartCheckout(&orga, 0)
	if err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	if err = Repo.CreateCheckoutSession(&models.CheckoutSession{
		SessionID: sessionID,
	}).Err(); err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	response = &models.StartSubscriptionResponse{
		SessionID: sessionID,
		Url:       url,
	}

done:
	return response, status, err
}

func GetPollSubscriptionSuccess(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (response router.Serde, status int, err error) {
	status = http.StatusOK
	var cs models.CheckoutSession
	var sessionID string = params.Get("sessionID").(string)

	if err = Repo.GetCheckoutSession(sessionID, &cs).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

	return response, status, err
}

func GetCheckoutSuccess(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	status := http.StatusOK
	sessionID := r.URL.Query().Get("session_id")
	msg := "Thank you for subscribing to Keystone!"

	err := repo.Transaction(func(Repo repo.IRepo) error {
		var cs models.CheckoutSession

		return Repo.
			GetCheckoutSession(sessionID, &cs).
			DeleteCheckoutSession(&cs).
			Err()
	})

	if err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			msg = "No such checkout session"
		} else {
			status = http.StatusInternalServerError
			msg = fmt.Sprintf("An error occurred : %s", err.Error())
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", strconv.Itoa(len(msg)))
	w.Write([]byte(msg))

	if status != http.StatusOK {
		w.WriteHeader(status)
	}
}

func GetCheckoutCancel(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	status := http.StatusOK
	sessionID := r.URL.Query().Get("session_id")
	msg := "You cancelled your subscriptio"

	err := repo.Transaction(func(Repo repo.IRepo) error {
		var cs models.CheckoutSession

		return Repo.
			GetCheckoutSession(sessionID, &cs).
			DeleteCheckoutSession(&cs).
			Err()
	})

	if err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			msg = "No such checkout session"
		} else {
			status = http.StatusInternalServerError
			msg = fmt.Sprintf("An error occurred : %s", err.Error())
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", strconv.Itoa(len(msg)))
	w.Write([]byte(msg))
	w.WriteHeader(status)
}

func PostStripeWebhook(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	var err error
	var event payment.Event

	p := payment.NewStripePayment()
	event, err = p.HandleEvent(r)
	fmt.Printf("event: %+v\n", event)

	switch event.Type {
	case payment.EventCheckoutComplete:
		organization := models.Organization{
			ID: event.OrganizationID,
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			Repo.GetOrganization(&organization)
			organization.SubscriptionID = string(event.SubscriptionID)

			Repo.UpdateOrganization(&organization)

			return Repo.Err()
		})

	case payment.EventSubscriptionPaid:
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			return Repo.
				GetOrganization(&organization).
				OrganizationSetPaid(&organization, true).
				Err()
		})

	case payment.EventSubscriptionUnpaid:
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			return Repo.
				GetOrganization(&organization).
				OrganizationSetPaid(&organization, false).
				Err()
		})

	case payment.EventSubscriptionCancelled:
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			return Repo.
				GetOrganization(&organization).
				OrganizationSetPaid(&organization, false).
				Err()
		})

	default:
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}
