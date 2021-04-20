// +build test

package nock

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
	"gopkg.in/h2non/gock.v1"
)

type GithubUserType struct {
	id        int
	UserEmail []*github.UserEmail
}

// Email      *string `json:"email,omitempty"`
// 	Primary    *bool   `json:"primary,omitempty"`
// 	Verified   *bool   `json:"verified,omitempty"`
// 	Visibility *string `json:"visibility,omitempty"`

func StartNock() {

	// defer gock.Off() // Flush pending mocks after test execution
	Email := "sfsdf@edfsf.com"
	Verified := true
	Visibility := "Visibility"

	UserEmail := github.UserEmail{
		Email:      &Email,
		Verified:   &Verified,
		Visibility: &Visibility}

	var emails []*github.UserEmail
	emails = append(emails, &UserEmail)

	gock.New("https://api.github.com").
		Persist().
		Get("/user").
		// Reply(200).
		ReplyFunc(func(resp *gock.Response) {
			fmt.Println("🚀 ~ file: start-nock_testing.go ~ line 48 ~ funcStartNock ~ resp", resp)

			resp.SetHeader("Content-Type", "application/json")
			resp.Status(200)

			objUser := &github.User{}

			userJSON := `
			{
  "login": "octocat",
  "id": 1,
  "node_id": "MDQ6VXNlcjE=",
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/octocat",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/octocat/followers",
  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
  "organizations_url": "https://api.github.com/users/octocat/orgs",
  "repos_url": "https://api.github.com/users/octocat/repos",
  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
  "received_events_url": "https://api.github.com/users/octocat/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa octocat",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "octocat@github.com",
  "hireable": false,
  "bio": "There once was...",
  "twitter_username": "monatheoctocat",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z",
  "private_gists": 81,
  "total_private_repos": 100,
  "owned_private_repos": 100,
  "disk_usage": 10000,
  "collaborators": 8,
  "two_factor_authentication": true,
  "plan": {
    "name": "Medium",
    "space": 400,
    "private_repos": 20,
    "collaborators": 0
  }
}
`

			json.Unmarshal(userJSON, &objUser)
			// fmt.Println("ksauth ~ start-nock_testing.go ~ objUser", objUser)

			userJSONString := strings.NewReader(userJSON)

			// json.NewDecoder(userJSONString).Decode(objUser)

			json.NewDecoder(userJSONString).Decode(objUser)
			fmt.Println("keystone ~ start-nock_testing.go ~ objUser", objUser)

			// user := new(github.User)

			// json.NewDecoder(userJSONString).Decode(user)

			// // login := "LOGIN"
			// // user.Login = &login

			// b, _ := json.Marshal(user)
			// fmt.Println("ksauth ~ start-nock_testing.go ~ b", string(b))

			// fmt.Println("ksauth ~ start-nock_testing.go ~ user", user)
			// fmt.Println("ksauth ~ start-nock_testing.go ~ user", user)
			// fmt.Println("ksauth ~ start-nock_testing.go ~ user", user.Login)

			resp.SetHeader("Content-Type", "application/json")
			resp.BodyString(userJSON)
			// resp.JSON(&{
			// 	id:        56883564,
			// 	UserEmail: emails})
		})
		// Reply(200).
		// JSON(&GithubUserType{
		// 	id:        56883564,
		// 	UserEmail: emails})

	gock.New("https://api.github.com").
		Persist().
		Get("/user/emails").
		ReplyFunc(func(resp *gock.Response) {
			fmt.Println("🚀 ~ file: start-nock_testing.go ~ line 48 ~ funcStartNock ~ resp", resp)

			resp.BodyString("looooooooooool")
			resp.JSON(&GithubUserType{
				id:        56883564,
				UserEmail: emails})
		})

	// Reply(200).
	// JSON(&GithubUserType{
	// 	id:        56883564,
	// 	UserEmail: emails})

}
