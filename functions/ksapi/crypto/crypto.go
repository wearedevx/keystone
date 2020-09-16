package crypto

import (
	. "strings"

	. "github.com/wearedevx/keystone/internal/models"

	"filippo.io/age"
	"filippo.io/age/agessh"
)

func EncryptForUser(user *User, payload []byte, out *[]byte) error {
	var recipient age.Recipient
	var err error

	if HasPrefix(user.Keys.Cipher, "ssh-") {
		recipient, err = agessh.ParseRecipient(user.Keys.Cipher)
	} else {
		recipient, err = age.ParseX25519Recipient(user.Keys.Cipher)
	}

	if err != nil {
		var sb Builder

		w, err := age.Encrypt(&sb, recipient)
		if err != nil {
			return err
		}

		w.Write(payload)
		w.Close()

		*out = []byte(sb.String())
	}

	return err
}
