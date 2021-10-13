package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	"github.com/wearedevx/keystone/api/internal/payment"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func PostSubscription(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (response router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "PostSubscription",
	}

	var sessionID, customerID, url string
	var seats int64

	p := payment.NewStripePayment()

	organizationName := params.Get("organizationName").(string)
	organization := models.Organization{
		Name: organizationName,
	}

	if err = Repo.
		GetOrganization(&organization).
		OrganizationCountMembers(&organization, &seats).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToGetResource(err)
		}

		goto done
	}

	if organization.User.UserID != user.UserID {
		status = http.StatusForbidden
		err = apierrors.ErrorPermissionDenied()
		goto done
	}

	if organization.CustomerID != "" && organization.SubscriptionID != "" && organization.Paid {
		status = http.StatusConflict
		err = apierrors.ErrorAlreadySubscribed()
		goto done
	}

	sessionID, customerID, url, err = p.StartCheckout(&organization, seats)
	if err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToStartCheckout(err)
		goto done
	}

	if err = Repo.
		OrganizationSetCustomer(&organization, customerID).
		CreateCheckoutSession(&models.CheckoutSession{
			SessionID: sessionID,
		}).Err(); err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToCreateResource(err)
		goto done
	}

	response = &models.StartSubscriptionResponse{
		SessionID: sessionID,
		Url:       url,
	}

done:
	return response, status, log.SetError(err)
}

func GetPollSubscriptionSuccess(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (response router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetPollSubscriptionSuccess",
	}

	var cs models.CheckoutSession
	var sessionID string = params.Get("sessionID").(string)

	if err = Repo.GetCheckoutSession(sessionID, &cs).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToGetResource(err)
		}
	}

	return response, status, log.SetError(err)
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
	msg := "You cancelled your subscription"

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

	switch event.Type {
	case payment.EventCheckoutComplete:
		organization := models.Organization{
			ID: event.OrganizationID,
		}
		session := models.CheckoutSession{}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			var seats int64
			var err error
			if err = Repo.
				GetCheckoutSession(event.SessionID, &session).
				Err(); err != nil {
				return err
			}

			session.Status = models.CheckoutSessionStatusSuccess

			if err = Repo.
				UpdateCheckoutSession(&session).
				OrganizationSetCustomer(&organization, string(event.CustomerID)).
				OrganizationSetSubscription(&organization, string(event.SubscriptionID)).
				OrganizationSetPaid(&organization, true).
				OrganizationCountMembers(&organization, &seats).
				Err(); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			err = apierrors.ErrorCheckoutCompleteFailed(err)
		}

	case payment.EventSubscriptionPaid:
		var seats int64
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			var err error

			if err = Repo.
				GetOrganization(&organization).
				OrganizationSetPaid(&organization, true).
				OrganizationCountMembers(&organization, &seats).
				Err(); err != nil {
				return err
			}

			if err = p.
				UpdateSubscription(event.SubscriptionID, seats); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			err = apierrors.ErrorSubscriptionPaidFailed(err)
		}

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

		if err != nil {
			err = apierrors.ErrorSubscriptionUnpaidFailed(err)
		}

	case payment.EventSubscriptionCanceled:
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			return Repo.
				GetOrganization(&organization).
				OrganizationSetPaid(&organization, false).
				Err()
		})

		if err != nil {
			err = apierrors.ErrorSubscriptionCanceledFailed(err)
		}

	default:
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func ManageSubscription(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	var url string
	var result models.ManageSubscriptionResponse

	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "ManageSubscription",
	}

	organizationName := params.Get("organizationName").(string)
	organization := models.Organization{
		Name: organizationName,
	}

	if err = Repo.
		GetOrganization(&organization).
		Err(); err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetResource(err)
		goto done
	}

	url, err = payment.
		NewStripePayment().
		GetManagementLink(&organization)
	if err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetManagementLink(err)
		goto done
	}

	result = models.ManageSubscriptionResponse{
		Url: url,
	}

done:
	return &result, status, log.SetError(err)
}
