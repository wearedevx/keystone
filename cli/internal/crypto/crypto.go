package crypto

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/cossacklabs/themis/gothemis/cell"
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

// Encrypts a file using a user-provided passphrase.
// `filepath` is the path to the file to be encrypted, and
// `passphrase` is the user-provided passphrase.
// It returns the encrypted content and an error.
func EncryptFile(filepath, passphrase string) (encrypted []byte, err error) {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return encrypted, fmt.Errorf("cannot read %s: %w", filepath, err)
	}

	data := base64.StdEncoding.EncodeToString(contents)

	scell, err := cell.SealWithPassphrase(passphrase)
	if err != nil {
		return encrypted, err
	}

	encrypted, err = scell.Encrypt([]byte(data), nil)
	if err != nil {
		return encrypted, err
	}

	return encrypted, nil
}

// Attempbs to decrypt a file using a user-porvided passphrase
// `filepath` is the path to the encrypted file
// `target` is the path to output the decrypted content
// `passphrase` is the user-provided passphrase
// It returns the decrypted content and an error
func DecryptFile(filepath, target, passphrase string) (err error) {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	scell, err := cell.SealWithPassphrase(passphrase)
	if err != nil {
		return err
	}

	decrypted, err := scell.Decrypt([]byte(contents), nil)
	if err != nil {
		return err
	}

	decrypted, err = base64.StdEncoding.DecodeString(string(decrypted))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(target, decrypted, 0600)

	return err
}
