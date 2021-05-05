package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wearedevx/keystone/api/pkg/repo"
	. "github.com/wearedevx/keystone/internal/models"
)

func EndScript() int {
	// os.Kill
	return 0
}

func GithubLoginSuccess() int {
	time.Sleep(3000 * time.Millisecond)

	lr := LoginRequest{}

	// Retrieve login_attemps in db
	Repo := new(repo.Repo)
	db := Repo.GetDb()

	if error := db.Last(&lr); error != nil {
		fmt.Println(error)
	}

	// simulate github POST on auth cloud function

	timeout := time.Duration(20 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", "http://localhost:9001/auth-redirect/?state="+lr.TemporaryCode+"&code=youpicode", nil)

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
