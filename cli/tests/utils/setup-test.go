package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/rogpeppe/go-internal/testscript"
	uuid "github.com/satori/go.uuid"
	"github.com/wearedevx/keystone/api/pkg/jwt"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetGcloudFuncAuthPidFilePath() string {
	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)

	return path.Join(os.TempDir(), "keystone_ksauth.pid")
}

func WaitAPIStart() error {
	started := waitForServerStarted("http://127.0.0.1:9001")

	if started {
		return nil
	}

	return errors.New("server start timeout")
}

func isServerResponse(serverURL string) bool {
	request, _ := http.NewRequest("GET", serverURL, nil)

	timeout := time.Duration(1 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	_, err := client.Do(request)

	// If it's started,

	return err == nil
}

func pollServer(serverURL string, c chan bool, maxAttempts int) {
	var done = false
	attemps := 0

	for !done {
		attemps = attemps + 1

		if attemps == maxAttempts {
			done = true
		}

		isServerStarted := isServerResponse(serverURL)

		if isServerStarted {
			done = true
			c <- true
		}

		time.Sleep(1 * time.Second)
	}
}

func waitForServerStarted(serverURL string) bool {
	const max_attempts int = 40
	var result bool

	c := make(chan bool)

	go pollServer(serverURL, c, max_attempts)

	result = <-c

	return result
}

func CreateFakeUserWithUsername(
	username string,
	accountType models.AccountType,
	env *testscript.Env,
) (err error) {
	Repo := new(repo.Repo)
	user := models.User{}

	if err = faker.FakeData(&user); err != nil {
		return err
	}

	keyPair, err := keys.New(keys.TypeEC)

	deviceUID := uuid.NewV4().String()
	device := "device-test-" + deviceUID
	user.Username = username
	user.AccountType = accountType
	user.UserID = fmt.Sprintf("%s@%s", user.Username, user.AccountType)
	user.Devices = []models.Device{
		{Name: device, UID: deviceUID, PublicKey: keyPair.Public.Value},
	}

	if err = Repo.GetOrCreateUser(&user).Err(); err != nil {
		fmt.Println(err)
		return err
	}

	token, err := jwt.MakeToken(user, deviceUID, time.Now())
	configDir := getConfigDir(env)
	pathToKeystoneFile := path.Join(configDir, "keystone2.yaml")

	pub := base64.StdEncoding.EncodeToString(keyPair.Public.Value)
	priv := base64.StdEncoding.EncodeToString(keyPair.Private.Value)

	/* #nosec */
	err = ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- fullname: `+user.Fullname+`
  account_type: "`+string(user.AccountType)+`"
  email: `+user.Email+`
  ext_id: "`+user.ExtID+`"
  fullname: `+user.Fullname+`
  user_id: `+user.UserID+`
  username: `+user.Username+`
auth_token: `+token+`
device: `+device+`
device_uid: `+deviceUID+`
public_key: !!binary `+pub+`
private_key: !!binary `+priv+`
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

	if err = faker.FakeData(&user); err != nil {
		return err
	}

	keyPair, err := keys.New(keys.TypeEC)

	device := "device-test"
	deviceUID := uuid.NewV4().String()
	user.ID = 0
	user.Email = "email@example.com"
	user.Devices = []models.Device{
		{Name: device, UID: deviceUID, PublicKey: keyPair.Public.Value},
		{
			Name:      "device-test-2",
			UID:       uuid.NewV4().String(),
			PublicKey: keyPair.Public.Value,
		},
	}

	if err := Repo.GetOrCreateUser(&user).Err(); err != nil {
		return err
	}

	for _, orga := range user.Organizations {
		orga.Paid = true
		Repo.GetDB().Save(&orga)
	}

	if err := Repo.Err(); err != nil {
		fmt.Println("Get Or Create User", err)
		os.Exit(1)
	}

	env.Setenv("USER_ID", user.UserID)
	fmt.Printf("user.UserID: %+v\n", user.UserID)

	token, err := jwt.MakeToken(user, deviceUID, time.Now())
	configDir := getConfigDir(env)
	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	pub := base64.StdEncoding.EncodeToString(keyPair.Public.Value)
	priv := base64.StdEncoding.EncodeToString(keyPair.Private.Value)

	/* #nosec */
	err = ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- fullname: `+user.Fullname+`
  account_type: "`+string(user.AccountType)+`"
  email: `+user.Email+`
  ext_id: "`+user.ExtID+`"
  fullname: `+user.Fullname+`
  user_id: `+user.UserID+`
  username: `+user.Username+`
auth_token: `+token+`
device: `+device+`
device_uid: `+deviceUID+`
public_key: !!binary `+pub+`
private_key: !!binary `+priv+`
current: 0
`), 0o660)

	if err != nil {
		fmt.Println("error writing accounts", err)
		return err
	}

	fmt.Println("Written", pathToKeystoneFile)

	return err
}

func CreateProject(env *testscript.Env) (err error) {
	Repo := new(repo.Repo)
	project := models.Project{}
	user := models.User{}

	if err = faker.FakeData(&project); err != nil {
		return err
	}

	if err = faker.FakeData(&user); err != nil {
		return err
	}

	keyPair, err := keys.New(keys.TypeEC)

	device := "device-test"
	deviceUID := uuid.NewV4().String()
	user.ID = 0
	user.Email = "email@example.com"
	user.Devices = []models.Device{
		{Name: device, UID: deviceUID, PublicKey: keyPair.Public.Value},
	}

	Repo.
		GetOrCreateUser(&user).
		GetOrCreateProject(&project).
		ProjectAddMembers(
			project,
			[]models.MemberRole{
				{MemberID: user.UserID, RoleID: 4},
			},
			user,
		)

	if err := Repo.Err(); err != nil {
		fmt.Println(
			"Get Or Create User, or Project, or add Member to Project",
			err,
		)
		os.Exit(1)
	}

	cwd := env.Getenv("WORK")

	pathToKeystoneFile := path.Join(cwd, "keystone.yaml")
	err = ioutil.WriteFile(
		pathToKeystoneFile,
		[]byte(`project_id: `+project.UUID+`
name: `+project.Name+`
env: []
files: []
ci_services: []`),
		0o660,
	)

	return nil
}

func getHomeDir(env *testscript.Env) string {
	return path.Join(env.Getenv("WORK"), "home")
}

func getConfigDir(env *testscript.Env) string {
	homeDir := getHomeDir(env)
	return path.Join(homeDir, ".config", "keystone")
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
	return os.MkdirAll(configDir, 0o770) // #nosec
}

func MakeOrgaFree(env *testscript.Env) error {
	userID := env.Getenv("USER_ID")
	Repo := new(repo.Repo)

	user := models.User{UserID: userID}
	if err := Repo.GetDB().Preload("Organizations").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	for _, orga := range user.Organizations {
		orga.Paid = false
		if err := Repo.GetDB().Save(&orga).Error; err != nil {
			return err
		}

	}

	return nil
}
