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
	router.GET("/healthcheck", GetHealthCheck)
	router.GET("/", AuthedHandler(GetUser))

	router.POST("/projects", AuthedHandler(PostProject))
	router.DELETE("/projects/:projectID/", AuthedHandler(DeleteProject))

	router.GET("/projects/:projectID/members", AuthedHandler(GetProjectsMembers))
	router.POST("/projects/:projectID/members", AuthedHandler(PostProjectsMembers))
	router.DELETE("/projects/:projectID/members", AuthedHandler(DeleteProjectsMembers))
	router.PUT("/projects/:projectID/members/role", AuthedHandler(PutMembersSetRole))
	router.GET("/projects/:projectID/environments", AuthedHandler(GetAccessibleEnvironments))

	router.GET("/environments/:envID/public-keys", AuthedHandler(GetEnvironmentPublicKeys))
	router.DELETE("/messages-expired", DeleteExpiredMessages)
	router.GET("/messages-will-expire", AlertMessagesWillExpire)
	router.POST("/messages", AuthedHandler(WriteMessages))

	router.GET("/roles", AuthedHandler(GetRoles))

	router.GET("/devices", AuthedHandler(GetDevices))
	router.DELETE("/devices/:name", AuthedHandler(DeleteDevice))

	router.POST("/login-request", PostLoginRequest)
	router.GET("/login-request", GetLoginRequest)
	router.GET("/auth-redirect/", GetAuthRedirect)
	router.POST("/complete", PostUserToken)

	router.POST("/users/exist", AuthedHandler(DoUsersExist))
	router.GET("/users/:userID/key", AuthedHandler(GetUserKeys))
	router.POST("/users/invite", AuthedHandler(PostInvite))

	router.GET("/projects/:projectID/messages/:device", AuthedHandler(GetMessagesFromProjectByUser))

	router.DELETE("/messages/:messageID", AuthedHandler(DeleteMessage))

	router.ServeHTTP(w, r)
}
