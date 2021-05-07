package utils

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/api/pkg/jwt"
	. "github.com/wearedevx/keystone/api/pkg/models"
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

func CreateAndLogUser(env *testscript.Env) error {
	Repo := new(repo.Repo)
	username := "LAbigael_" + uuid.NewV4().String()
	userID := username + "@github"

	ExtID := "56883564" + uuid.NewV4().String()
	Username := "LAbigael_" + uuid.NewV4().String()

	var user1 *User = &User{
		ExtID:       ExtID,
		UserID:      userID,
		AccountType: "github",
		Username:    Username,
		Fullname:    "Test user1",
		Email:       "abigael.laldji@protonmail.com",
	}

	Repo.GetOrCreateUser(user1)

	token, _ := MakeToken(*user1)

	configDir := getConfigDir(env)

	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	err := ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- Fullname: `+user1.Fullname+`
  account_type: "`+string(user1.AccountType)+`"
  email: `+user1.Email+`
  ext_id: "`+ExtID+`"
  fullname: `+user1.Fullname+`
  user_id: `+user1.UserID+`
  username: `+Username+`
auth_token: `+token+`
current: 0
`), 0o777)

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

	// Create config folder
	return os.MkdirAll(configDir, 0777)
}
