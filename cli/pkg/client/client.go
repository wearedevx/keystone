package client

var APIURL string //= "http://localhost:9001"

type KeystoneClient interface {
	Project(projectID string) *Project
	Roles() *Roles
	Users() *Users
	Messages() *Messages
	Devices() *Devices
	Organizations() *Organizations
}
