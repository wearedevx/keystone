package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

	organizationName := params.Get("organizationName")
	var organization models.Organization

	if organizationName == "" {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(nil)
		goto done
	}

	organization = models.Organization{
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

	if organization.CustomerID != "" && organization.SubscriptionID != "" &&
		organization.Paid {
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
		URL:       url,
	}

done:
	return response, status, log.SetError(err)
}

func GetPollSubscriptionSuccess(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (response router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetPollSubscriptionSuccess",
	}

	var cs models.CheckoutSession
	sessionID := params.Get("sessionID")

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
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	status = http.StatusOK

	sessionID := r.URL.Query().Get("session_id")
	msg := "Thank you for subscribing to Keystone!"

	var cs models.CheckoutSession

	if err = Repo.
		GetCheckoutSession(sessionID, &cs).
		DeleteCheckoutSession(&cs).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			msg = "No such checkout session"
		} else {
			status = http.StatusInternalServerError
			msg = fmt.Sprintf("An error occurred: %s", err.Error())
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", strconv.Itoa(len(msg)))
	w.Write([]byte(msg))

	return status, err
}

func GetCheckoutCancel(
	w http.ResponseWriter,
	r *http.Request,
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	status = http.StatusOK
	sessionID := r.URL.Query().Get("session_id")
	msg := "You cancelled your subscription"

	var cs models.CheckoutSession

	if err = Repo.
		GetCheckoutSession(sessionID, &cs).
		DeleteCheckoutSession(&cs).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			msg = "No such checkout session"
		} else {
			status = http.StatusInternalServerError
			msg = fmt.Sprintf("An error occurred: %s", err.Error())
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", strconv.Itoa(len(msg)))
	w.Write([]byte(msg))

	return status, err
}

func PostStripeWebhook(
	w http.ResponseWriter,
	r *http.Request,
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	status = http.StatusOK
	var event payment.Event

	p := payment.NewStripePayment()
	event, err = p.HandleEvent(r)

	if err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	switch event.Type {
	case payment.EventCheckoutComplete:
		organization := models.Organization{
			ID: event.OrganizationID,
		}
		session := models.CheckoutSession{}

		var seats int64
		if err = Repo.
			GetCheckoutSession(event.SessionID, &session).
			Err(); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorCheckoutCompleteFailed(err)
			goto done
		}

		session.Status = models.CheckoutSessionStatusSuccess

		if err = Repo.
			UpdateCheckoutSession(&session).
			OrganizationSetCustomer(&organization, string(event.CustomerID)).
			OrganizationSetSubscription(&organization, string(event.SubscriptionID)).
			OrganizationSetPaid(&organization, true).
			OrganizationCountMembers(&organization, &seats).
			Err(); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorCheckoutCompleteFailed(err)
		}

	case payment.EventSubscriptionPaid:
		var seats int64
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		if err = Repo.
			GetOrganization(&organization).
			OrganizationSetPaid(&organization, true).
			OrganizationCountMembers(&organization, &seats).
			Err(); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorSubscriptionPaidFailed(err)
			goto done
		}

		if err = p.
			UpdateSubscription(event.SubscriptionID, seats); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorSubscriptionPaidFailed(err)
		}

	case payment.EventSubscriptionUnpaid:
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		if err = Repo.
			GetOrganization(&organization).
			OrganizationSetPaid(&organization, false).
			Err(); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorSubscriptionUnpaidFailed(err)
		}

	case payment.EventSubscriptionCanceled:
		organization := models.Organization{
			SubscriptionID: string(event.SubscriptionID),
		}

		if err = Repo.
			GetOrganization(&organization).
			OrganizationSetPaid(&organization, false).
			Err(); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorSubscriptionCanceledFailed(err)
		}

	default:
	}

done:
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	return status, err
}

func ManageSubscription(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	var url string
	var result *models.ManageSubscriptionResponse

	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "ManageSubscription",
	}

	organizationName := params.Get("organizationName")
	organization := models.Organization{
		Name: organizationName,
	}

	if err = Repo.
		GetOrganization(&organization).
		Err(); err != nil {
		if errors.Is(repo.ErrorNotFound, err) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToGetResource(err)
		}
		goto done
	}

	if organization.UserID != user.ID {
		status = http.StatusForbidden
		err = apierrors.ErrorPermissionDenied()
		goto done
	}

	if !organization.Paid {
		status = http.StatusForbidden
		err = apierrors.ErrorNeedsUpgrade()
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

	result = &models.ManageSubscriptionResponse{
		URL: url,
	}

done:
	return result, status, log.SetError(err)
}
