package client

type KeystoneClientImpl struct {
	r requester
}

func NewKeystoneClient(userID string, jwtToken string) KeystoneClient {
	return &KeystoneClientImpl{
		r: newRequester(userID, jwtToken),
	}
}

func (client *KeystoneClientImpl) Roles() *Roles {
	return &Roles{
		r: client.r,
	}
}

func (client *KeystoneClientImpl) Project(projectId string) *Project {
	return &Project{
		r:  client.r,
		id: projectId,
	}
}

func (client *KeystoneClientImpl) Users() *Users {
	return &Users{
		r: client.r,
	}
}

func (client *KeystoneClientImpl) Messages() *Messages {
	return &Messages{
		r: client.r,
	}
}
