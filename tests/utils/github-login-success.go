package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
)

func EndScript() int {
	// os.Kill
	return 0
}

func GithubLoginSuccess() int {

	fmt.Println(" keystone ~ github-login-success.go ~  start")

	time.Sleep(3000 * time.Millisecond)
	fmt.Println(" keystone ~ github-login-success.go ~  os.Getpid() !", os.Getpid())

	lr := LoginRequest{}

	// Retrieve login_attemps in db
	Repo := new(repo.Repo)
	Repo.Connect()
	db := Repo.GetDb()

	if error := db.Last(&lr); error != nil {
		fmt.Println(error)
	}

	fmt.Println("🚀 ~ file: github-login-success.go ~ line 19 ~ funcGithubLoginSuccess ~ loginRequest", lr)

	// simulate github POST on auth cloud function

	timeout := time.Duration(20 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", "http://localhost:9000/auth-redirect/"+lr.TemporaryCode+"/?code=youpicode", nil)

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusOK {
	}

	fmt.Println("github login success End")

	return 0
}
