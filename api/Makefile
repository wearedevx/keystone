ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= ks
DB_PASSWORD ?= ks
DB_NAME ?= ks
DB_DIALECT ?= postgres
DOMAIN ?= localhost:9001

STRIPE_KEY ?= sk_test

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
DOMAIN_FLAG := -X $(CONSTANTS_PKG).Domain=$(DOMAIN)

REDIS_HOST_FLAG :=  -X $(REDIS_PKG).redisHost=$(REDIS_HOST)
REDIS_PORT_FLAG :=  -X $(REDIS_PKG).redisPort=$(REDIS_PORT)
REDIS_INDEX_FLAG := -X $(REDIS_PKG).redisIndex=$(REDIS_INDEX)

STRIPE_KEY_FLAG := -X $(PAYMENT_PKG).stripeKey=$(STRIPE_KEY) 
STRIPE_WEBHOOK_FLAG := -X $(PAYMENT_PKG).stripeWebhookSecret=$(STRIPE_WEBHOOK_SECRET)
STRIPE_PRICE_FLAG := -X $(PAYMENT_PKG).stripeSubscriptionPrice=$(STRIPE_PRICE)
 
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
	$(STRIPE_PRICE_FLAG) \
	$(STRIPE_WEBHOOK_FLAG) \
	$(KS_TTL_FLAG) \
	$(DOMAIN_FLAG) 

build:
	mkdir ./build
	go clean
	go generate
	go build -ldflags "$(LDFLAGS)" -v -o ./build/server

build-debug:
	go generate
	go build -gcflags="all=-N -l" -ldflags="$(LDFLAGS)" -o ./server
	chmod +x ./server

install:
	chmod +x server
	cp ks $(PREFIX)/bin

run-test:
	go run -tags test -ldflags "$(LDFLAGS)" main.go

test:
	if [ -f "${TMPDIR}keystone_gorm.db" ]; then rm "${TMPDIR}keystone_gorm.db"; fi
	touch "${TMPDIR}keystone_gorm.db"
	go test -tags test -coverprofile cover.out.tmp ./...
	cat cover.out.tmp | grep -v "generated_" | grep -v "_testing" > cover.out
	rm cover.out.tmp

run:
	reflex \
		-R generated \
		-s -- \
		sh -c 'go generate && go run -ldflags "$(LDFLAGS)" main.go'

run-debug:
	reflex \
		-R server \
		-R generated \
		-s -- \
		sh -c "make build-debug && dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient --log exec --continue ./server"
