package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetHealthCheck(
	w http.ResponseWriter,
	_ *http.Request,
	_ httprouter.Params,
) {
	http.Error(w, "", http.StatusOK)
}
