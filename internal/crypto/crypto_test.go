// Package crypto provides ...
package crypto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wearedevx/keystone/internal/models"
)

func TestEncryptForUser(t *testing.T) {
	var user models.User = models.User{
		// Keys: models.KeyRing{
		// 	Sign:   "ssh-rsa ",
		// 	Cipher: "ssh-rsa ",
		// },
	}

	inString := "Hello user"
	in := bytes.NewBufferString(inString)
	fmt.Printf("There are %d bytes to encrypt\n", in.Len())
	var out bytes.Buffer

	n, err := EncryptForUser(&user, in, &out)
	fmt.Printf("Encrypted %d bytes\n", n)
	fmt.Printf("Error ? %+v\n", err)

	fmt.Println(out)
}
