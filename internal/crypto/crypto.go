package crypto

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	. "strings"

	"github.com/mitchellh/go-homedir"
	. "github.com/wearedevx/keystone/internal/models"

	"filippo.io/age"
	"filippo.io/age/agessh"
)

func EncryptForUser(user *User, payload io.Reader, out io.Writer) (int64, error) {
	return EncryptForPublicKey(string(user.PublicKey), payload, out)
}

func EncryptForPublicKey(publicKey string, payload io.Reader, out io.Writer) (int64, error) {
	var recipient age.Recipient
	var err error
	var n int64 = 0

	if HasPrefix(publicKey, "ssh-") {
		recipient, err = agessh.ParseRecipient(publicKey)
	} else {
		recipient, err = age.ParseX25519Recipient(publicKey)
	}

	if err == nil {
		w, err := age.Encrypt(out, recipient)
		if err != nil {

			return n, err
		}

		n, err = io.Copy(w, payload)

		if err != nil {
			return n, err
		}

		err = w.Close()
		if err != nil {
			return n, err
		}
	}

	return n, err
}

func saveTempPublicKey(publicKey string) (string, error) {
	osTempDir := os.TempDir()

	p, err := ioutil.TempFile(osTempDir, "keystone.*.pub")
	if err != nil {
		return "", err
	}

	defer p.Close()

	// filePath := path.Join(osTempDir, p.Name())
	p.Write([]byte(publicKey))

	return p.Name(), nil
}

func extractPublicKeyFromFile(filepath string) (string, error) {
	cmd := exec.Command("ssh-keygen", "-yef", filepath)

	var out bytes.Buffer
	var serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr

	err := cmd.Run()

	if err != nil {
		fmt.Println(serr.String())
	}

	return out.String(), err
}

func findPrivateKey(publicKey string) []byte {
	pkPath, err := saveTempPublicKey(publicKey)
	if err != nil {
		fmt.Printf("Failed to save temporary public key at `%s`\n", pkPath)
		panic(err)
	}
	// defer os.RemoveAll(pkPath)

	pk, err := extractPublicKeyFromFile(pkPath)
	if err != nil {
		fmt.Printf("Failed to extract public key from file at `%s`\n", pkPath)
		fmt.Println(pk)
		panic(err)
	}

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("No homedir!")
		panic(err)
	}

	sshDirPath := path.Join(home, ".ssh")
	files, err := ioutil.ReadDir(sshDirPath)
	if err != nil {
		fmt.Printf("Cannot read directory: %s\n", sshDirPath) //
		panic(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			filePath := path.Join(sshDirPath, f.Name())
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Failed opening file at `%s`\n", filePath)
				panic(err)
			}

			defer file.Close()

			reader := bufio.NewReader(file)
			line, _, err := reader.ReadLine()

			if err != nil {
				fmt.Println("Cannot read line")
				panic(err)
			}

			if strings.Contains(string(line), "PRIVATE KEY") {
				pkCandidate, _ := extractPublicKeyFromFile(filePath)

				if pkCandidate == pk {
					content, err := ioutil.ReadFile(filePath)
					if err != nil {
						fmt.Printf("Failed to read file content at `%s`\n", filePath)
						panic(err)
					}
					fmt.Printf("Found private key at `%s`\n", filePath)

					return content
				}
			}
		}
	}

	return []byte("")
}

func DecryptWithPublicKey(publicKey []byte, payload []byte, out interface{}) error {
	var identity age.Identity
	var err error

	pk := string(publicKey)

	if HasPrefix(pk, "ssh-") {
		identity, err = agessh.ParseIdentity(findPrivateKey(pk))
	} else {
		identity, err = age.ParseX25519Identity(string(findPrivateKey(pk)))
	}

	if err != nil {
		return err
	}

	buf := bytes.NewReader(payload)
	output := &bytes.Buffer{}

	r, err := age.Decrypt(buf, identity)

	if err != nil {
		return err
	}

	if _, err := io.Copy(output, r); err != nil {
		log.Fatalf("Failed to read encrypted file: %v", err)
		panic(err)
		// return err
	}

	return json.NewDecoder(output).Decode(out)
}
