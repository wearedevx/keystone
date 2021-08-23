package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/api/pkg/jwt"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetGcloudFuncAuthPidFilePath() string {
	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)

	return path.Join(os.TempDir(), "keystone_ksauth.pid")
}

func WaitAPIStart() {
	waitForServerStarted("http://127.0.0.1:9001")
}

func isServerResponse(serverUrl string) bool {
	request, _ := http.NewRequest("GET", serverUrl, nil)

	timeout := time.Duration(1 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	_, err := client.Do(request)

	// If it's started,

	return err == nil
}

func pollServer(serverUrl string, c chan bool, maxAttempts int) {
	var done bool = false
	attemps := 0

	for !done {
		attemps = attemps + 1

		if attemps == maxAttempts {
			done = true
		}

		isServerStarted := isServerResponse(serverUrl)

		if isServerStarted {
			done = true
			c <- true
		}

		time.Sleep(1 * time.Second)
	}
}

func waitForServerStarted(serverUrl string) {
	const max_attempts int = 20

	c := make(chan bool)

	go pollServer(serverUrl, c, max_attempts)

	// result := true
	<-c
}

func CreateFakeUserWithUsername(username string, accountType models.AccountType, env *testscript.Env) (err error) {
	Repo := new(repo.Repo)
	user := models.User{}

	faker.FakeData(&user)
	keyPair, err := keys.New(keys.TypeEC)

	user.Username = username
	user.AccountType = accountType
	user.UserID = fmt.Sprintf("%s@%s", user.Username, user.AccountType)
	user.PublicKey = keyPair.Public.Value

	if err = Repo.GetOrCreateUser(&user).Err(); err != nil {
		return err
	}

	token, err := jwt.MakeToken(user)
	configDir := getConfigDir(env)
	pathToKeystoneFile := path.Join(configDir, "keystone2.yaml")

	pub := base64.StdEncoding.EncodeToString(keyPair.Public.Value)
	priv := base64.StdEncoding.EncodeToString(keyPair.Private.Value)

	err = ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- fullname: `+user.Fullname+`
  account_type: "`+string(user.AccountType)+`"
  email: `+user.Email+`
  ext_id: "`+user.ExtID+`"
  fullname: `+user.Fullname+`
  user_id: `+user.UserID+`
  username: `+user.Username+`
  public_key: !!binary `+pub+`
  private_key: !!binary `+priv+`
auth_token: `+token+`
current: 0
`), 0o666)

	if err != nil {
		fmt.Println("error wrinting user account", err)

		return err
	}

	fmt.Println("Written", pathToKeystoneFile)

	return nil
}

func CreateAndLogUser(env *testscript.Env) (err error) {
	Repo := new(repo.Repo)
	user := models.User{}

	faker.FakeData(&user)
	keyPair, err := keys.New(keys.TypeEC)

	user.ID = 0
	user.Email = "email@example.com"
	user.PublicKey = keyPair.Public.Value

	Repo.GetOrCreateUser(&user)

	if err := Repo.Err(); err != nil {
		fmt.Println("Get Or Create User", err)
		os.Exit(1)
	}

	env.Setenv("USER_ID", user.UserID)

	token, err := jwt.MakeToken(user)
	configDir := getConfigDir(env)
	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	pub := base64.StdEncoding.EncodeToString(keyPair.Public.Value)
	priv := base64.StdEncoding.EncodeToString(keyPair.Private.Value)

	err = ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- fullname: `+user.Fullname+`
  account_type: "`+string(user.AccountType)+`"
  email: `+user.Email+`
  ext_id: "`+user.ExtID+`"
  fullname: `+user.Fullname+`
  user_id: `+user.UserID+`
  username: `+user.Username+`
  public_key: !!binary `+pub+`
  private_key: !!binary `+priv+`
auth_token: `+token+`
current: 0
`), 0o660)

	if err != nil {
		fmt.Println("error writing accounts", err)
		return err
	}

	fmt.Println("Written", pathToKeystoneFile)

	return err
}

func getHomeDir(env *testscript.Env) string {
	return path.Join(env.Getenv("WORK"), "home")
}

func getConfigDir(env *testscript.Env) string {
	homeDir := getHomeDir(env)
	return path.Join(homeDir, ".config")
}

func SetupEnvVars(env *testscript.Env) error {
	homeDir := getHomeDir(env)
	configDir := getConfigDir(env)
	osTmpDir := os.TempDir()

	// Set home dir for test
	env.Setenv("GOPATH", "/DFDFDF")
	env.Setenv("HOME", homeDir)
	env.Setenv("DB_PORT", os.Getenv("DB_PORT"))
	env.Setenv("TMPDIR", osTmpDir)
	env.Setenv("GCLOUDFILE", GetGcloudFuncAuthPidFilePath())
	env.Setenv("NOSPIN", "true")
	env.Setenv("KSCOLORS", "off")

	// Create config folder
	return os.MkdirAll(configDir, 0o770)
}
