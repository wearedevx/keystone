package crypto

import (
	"fmt"

	"github.com/cossacklabs/themis/gothemis/errors"
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

func EncryptMessage(senderPrivateKey []byte, recipientPublicKey []byte, payload []byte) (msg []byte, err error) {
	private := keys.PrivateKey{Value: senderPrivateKey}
	fmt.Printf("private: %+v\n", private)
	public := keys.PublicKey{Value: recipientPublicKey}
	fmt.Printf("public: %+v\n", public)

	secureMessage := message.New(&private, &public)

	fmt.Printf("payload: %+v\n", payload)
	p, err := secureMessage.Wrap(payload)
	fmt.Printf("p: %+v\n", p)
	fmt.Printf("err: %+v\n", err)
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
