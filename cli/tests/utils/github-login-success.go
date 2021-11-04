package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// EndScript function
func EndScript() int {
	// os.Kill
	return 0
}

// GithubLoginSuccess function fakes a successfull login on github
func GithubLoginSuccess() int {
	time.Sleep(3000 * time.Millisecond)

	lrs := []models.LoginRequest{}

	// Retrieve login_attemps in db
	Repo := new(repo.Repo)
	db := Repo.GetDb()

	if error := db.Find(&lrs).Error; error != nil {
		fmt.Println("Error :", error)
	}

	// simulate github POST on auth cloud function
	for _, lr := range lrs {
		timeout := time.Duration(20 * time.Second)

		client := http.Client{
			Timeout: timeout,
		}
		state := models.AuthState{
			TemporaryCode: lr.TemporaryCode,
			Version:       "test",
		}

		codedState, err := state.Encode()
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequest(
			"GET",
			"http://localhost:9001/auth-redirect/?state="+codedState+"&code=youpicode",
			nil,
		)

		if err == nil {
			fmt.Println("reguest auth-redirect")
			_, err = client.Do(request)
		}

		if err != nil {
			panic(err)
		}

		fmt.Println("github login success End")
	}

	return 0
}
