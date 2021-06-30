ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif


BASE := github.com/wearedevx/keystone/cli
CLIENT_PKG := $(BASE)/pkg/client
AUTH_PKG := $(BASE)/pkg/client/auth

APIFLAG := -X $(CLIENT_PKG).ApiURL=$KSAPI_URL
GITHUB_CLIENT_ID_FLAG := -X $(PKG_AUTH).githubClientId=$GITHUB_CLIENT_ID
GITHUB_CLIENT_SECRET_FLAG := -X $(AUTH_PKG).githubClientSecret=$GITHUB_CILENT_SECRET
GITLAB_CLIENT_ID_FLAG := -X $(AUTH_PKG).gitlabClientId=$GITLAB_CLIENT_ID
GITLAB_CLIENT_SECRET_FLAG := -X $(AUTH_PKG).gitlabClientSecret=$GITLAB_CILENT_SECRET

LDFLAGS := $(APIFLAG) $(GITHUB_CLIENT_ID_FLAG) $(GITHUB_CLIENT_SECRET_FLAG) $(GITLAB_CLIENT_ID_FLAG) $(GITLAB_CILENT_SECRET_FLAG)

build:
	go clean
	go mod edit -replace=github.com/wearedevx/keystone/cmd=./cmd
	echo "$(LDFLAGS)"
	go build -ldflags "$(LDFLAGS)" -o ks

install:
	chmod +x ks
	cp ks $(PREFIX)/bin

run:
	go run main.go
