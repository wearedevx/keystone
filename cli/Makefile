ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

build:
	go clean
	go mod edit -replace=github.com/wearedevx/keystone/cmd=./cmd
	echo "$(LDFLAGS)"
	go build -ldflags "$(LDFLAGS)" -o ks
	chmod +x ks

install:
	cp ks $(PREFIX)/bin

run:
	go run main.go
