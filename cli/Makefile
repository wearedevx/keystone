ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

build:
	go clean
	go mod edit -replace=github.com/wearedevx/keystone/cmd=./cmd
	echo "$(KSLDFLAGS)"
	go build -ldflags "$(KSLDFLAGS)" -o ks

install:
	chmod +x ks
	cp ks $(PREFIX)/bin

run:
	go run main.go
