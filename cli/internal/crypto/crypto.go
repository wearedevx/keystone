package crypto

import (
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

func EncryptMessage(senderPrivateKey []byte, recipientPublicKey []byte, payload string) (msg []byte, err error) {
	private := keys.PrivateKey{Value: senderPrivateKey}
	public := keys.PublicKey{Value: recipientPublicKey}

	secureMessage := message.New(&private, &public)

	p, err := secureMessage.Wrap([]byte(payload))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func DecrypteMessage(recipientPrivateKey []byte, senderPublicKey []byte, msg []byte) (payload string, err error) {
	private := keys.PrivateKey{Value: recipientPrivateKey}
	public := keys.PublicKey{Value: senderPublicKey}

	secureMessage := message.New(&private, &public)

	p, err := secureMessage.Unwrap(msg)
	if err != nil {
		return "", err
	}

	return string(p), nil
}
