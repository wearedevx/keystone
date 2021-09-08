ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif


BASE := github.com/wearedevx/keystone/api
REPO_PKG := $(BASE)/pkg/repo
JWT_PKG := $(BASE)/pkg/jwt
CONSTANTS_PKG := $(BASE)/internal/constants

DB_HOST_FLAG := -X $(REPO_PKG).dbHost=$(DB_HOST)
DB_PORT_FLAG := -X $(REPO_PKG).dbPort=$(DB_PORT)
DB_USER_FLAG := -X $(REPO_PKG).dbUser=$(DB_USER)
DB_PASSWORD_FLAG := -X $(REPO_PKG).dbPassword=$(DB_PASSWORD)
DB_NAME_FLAG := -X $(REPO_PKG).dbName=$(DB_NAME)
DB_DIALECT_FLAG := -X $(REPO_PKG).dbDialect=$(DB_DIALECT)
DB_DRIVERNAME_FLAG := -X $(REPO_PKG).dbDriverName=$(DB_DRIVERNAME)

JWT_SALT_FLAG := -X $(JWT_PKG).salt=$(JWT_SALT)

KS_TTL_FLAG := -X $(CONSTANTS_PKG).KsTTLHeader=$(X_KS_TTL)
 
LDFLAGS := $(DB_HOST_FLAG) \
	$(DB_PORT_FLAG) \
	$(DB_USER_FLAG) \
	$(DB_PASSWORD_FLAG) \
	$(DB_NAME_FLAG) \
	$(DB_DIALECT_FLAG) \
	$(DB_DRIVERNAME_FLAG) \
	$(JWT_SALT_FLAG) \
	$(KS_TTL_FLAG)

build:
	go clean
	go build -ldflags "$(LDFLAGS)" -o server

install:
	chmod +x server
	cp ks $(PREFIX)/bin

run:
	go run -ldflags "$(LDFLAGS)" main.go
