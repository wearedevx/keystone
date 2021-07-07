FROM cossacklabs/ci-py-go-themis

RUN go get -u github.com/mitranim/gow
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN sudo cp migrate.linux-amd64 /usr/local/bin/migrate
