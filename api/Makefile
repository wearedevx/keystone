ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= ks
DB_PASSWORD ?= ks
DB_NAME ?= ks
DB_DIALECT ?= postgres

JWT_SALT ?= aaP|**P1n}1tqW

REDIS_HOST ?= localhost
REDIS_PORT ?= 6379
REDIS_INDEX ?= 0

BASE := github.com/wearedevx/keystone/api
REPO_PKG := $(BASE)/pkg/repo
JWT_PKG := $(BASE)/pkg/jwt
CONSTANTS_PKG := $(BASE)/internal/constants
REDIS_PKG := $(BASE)/internal/redis
PAYMENT_PKG := $(BASE)/internal/payment

DB_HOST_FLAG := -X $(REPO_PKG).dbHost=$(DB_HOST)
DB_PORT_FLAG := -X $(REPO_PKG).dbPort=$(DB_PORT)
DB_USER_FLAG := -X $(REPO_PKG).dbUser=$(DB_USER)
DB_PASSWORD_FLAG := -X $(REPO_PKG).dbPassword=$(DB_PASSWORD)
DB_NAME_FLAG := -X $(REPO_PKG).dbName=$(DB_NAME)
DB_DIALECT_FLAG := -X $(REPO_PKG).dbDialect=$(DB_DIALECT)
DB_DRIVERNAME_FLAG := -X $(REPO_PKG).dbDriverName=$(DB_DRIVERNAME)

JWT_SALT_FLAG := -X $(JWT_PKG).salt=$(JWT_SALT)

KS_TTL_FLAG := -X $(CONSTANTS_PKG).KsTTLHeader=$(X_KS_TTL)

REDIS_HOST_FLAG :=  -X $(REDIS_PKG).redisHost=$(REDIS_HOST)
REDIS_PORT_FLAG :=  -X $(REDIS_PKG).redisPort=$(REDIS_PORT)
REDIS_INDEX_FLAG := -X $(REDIS_PKG).redisIndex=$(REDIS_INDEX)

STRIPE_KEY_FLAG := -X $(PAYMENT_PKG).stripeKey=$(STRIPE_KEY) 
STRIPE_WEBHOOK_FLAG := -X $(PAYMENT_PKG).stripeWebhookSecret=$(STRIPE_WEBHOOK_SECRET)
 
LDFLAGS := $(DB_HOST_FLAG) \
	$(DB_PORT_FLAG) \
	$(DB_USER_FLAG) \
	$(DB_PASSWORD_FLAG) \
	$(DB_NAME_FLAG) \
	$(DB_DIALECT_FLAG) \
	$(DB_DRIVERNAME_FLAG) \
	$(JWT_SALT_FLAG) \
	$(REDIS_HOST_FLAG) \
	$(REDIS_PORT_FLAG) \
	$(REDIS_INDEX_FLAG) \
	$(STRIPE_KEY_FLAG) \
	$(STRIPE_WEBHOOK_FLAG) \
	$(KS_TTL_FLAG)

build:
	mkdir ./build
	go clean
	go build -ldflags "$(LDFLAGS)" -v -o ./build/server

build-debug:
	go clean
	go build -gcflags="all=-N -l" -ldflags="$(LDFLAGS)" -v -o ./server

install:
	chmod +x server
	cp ks $(PREFIX)/bin

run-test:
	go run -tags test -ldflags "$(LDFLAGS)" main.go

run:
	go run -ldflags "$(LDFLAGS)" main.go

run-debug:
	dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec ./server
