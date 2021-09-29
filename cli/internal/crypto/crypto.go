package crypto

import (
	"github.com/cossacklabs/themis/gothemis/errors"
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

func EncryptMessage(senderPrivateKey []byte, recipientPublicKey []byte, payload []byte) (msg []byte, err error) {
	private := keys.PrivateKey{Value: senderPrivateKey}
	public := keys.PublicKey{Value: recipientPublicKey}

	secureMessage := message.New(&private, &public)

	p, err := secureMessage.Wrap(payload)
	if err != nil {
		return nil, err
	}

	if len(p) == 0 {
		return nil, errors.New("crypto: invalid message length")
	}

	return p, nil
}

func DecryptMessage(recipientPrivateKey []byte, senderPublicKey []byte, msg []byte) (payload []byte, err error) {
	private := keys.PrivateKey{Value: recipientPrivateKey}
	public := keys.PublicKey{Value: senderPublicKey}

	secureMessage := message.New(&private, &public)

	p, err := secureMessage.Unwrap(msg)
	if err != nil {
		return []byte(""), err
	}

	return p, nil
}
