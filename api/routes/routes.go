// Package p contains an HTTP Cloud Function.
package routes

import (
	"net/http"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/julienschmidt/httprouter"

	. "github.com/wearedevx/keystone/api/controllers"
	. "github.com/wearedevx/keystone/api/internal/router"
)

// Auth shows the code to copy paste into the cli
func CreateRoutes(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()

	router.POST("/", PostUser)
	router.GET("/", AuthedHandler(GetUser))

	router.POST("/projects", AuthedHandler(PostProject))

	router.GET("/projects/:projectID/public-keys", AuthedHandler(GetProjectsPublicKeys))
	router.GET("/projects/:projectID/members", AuthedHandler(GetProjectsMembers))
	router.POST("/projects/:projectID/members", AuthedHandler(PostProjectsMembers))
	router.DELETE("/projects/:projectID/members", AuthedHandler(DeleteProjectsMembers))
	router.PUT("/projects/:projectID/members/role", AuthedHandler(PutMembersSetRole))

	router.GET("/environments/:envID/public-keys", AuthedHandler(GetEnvironmentPublicKeys))
	router.POST("/environments/:envID/messages", AuthedHandler(WriteMessages))

	router.GET("/roles", AuthedHandler(GetRoles))

	router.POST("/login-request", PostLoginRequest)
	router.GET("/login-request", GetLoginRequest)
	router.GET("/auth-redirect/", GetAuthRedirect)
	router.POST("/complete", PostUserToken)

	router.POST("/users/exist", AuthedHandler(DoUsersExist))

	router.GET("/projects/:projectID/messages", AuthedHandler(GetMessagesFromProjectByUser))

	router.DELETE("/messages/:messageID", AuthedHandler(DeleteMessage))

	router.ServeHTTP(w, r)
}
