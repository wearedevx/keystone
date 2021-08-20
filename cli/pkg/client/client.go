package client

var ApiURL string //= "http://localhost:9001"

type KeystoneClient interface {
	Project(projectId string) *Project
	Roles() *Roles
	Users() *Users
	Messages() *Messages
	Devices() *Devices
}
