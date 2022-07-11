# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.16-buster as builder

ARG GOOGLE_APPLICATION_CREDENTIALS
ARG DB_HOST
ARG DB_PORT
ARG DB_NAME
ARG DB_USER
ARG DB_PASSWORD
ARG JWT_SALT
ARG X_KS_TTL
ARG REDIS_HOST
ARG REDIS_PORT
ARG REDIS_INDEX
ARG STRIPE_KEY
ARG STRIPE_WEBHOOK_SECRET
ARG STRIPE_PRICE
ARG X_KS_TTL

ENV GOOGLE_APPLICATION_CREDENTIALS=$GOOGLE_APPLICATION_CREDENTIALS \
    DB_HOST=$DB_HOST \
    DB_PORT=$DB_PORT \
    DB_NAME=$DB_NAME \
    DB_USER=$DB_USER \
    DB_PASSWORD=$DB_PASSWORD \
    JWT_SALT=$JWT_SALT \
    X_KS_TTL=$X_KS_TTL \
    REDIS_HOST=$REDIS_HOST \
    REDIS_PORT=$REDIS_PORT \
    REDIS_INDEX=$REDIS_INDEX \
    STRIPE_KEY=$STRIPE_KEY \
    STRIPE_WEBHOOK_SECRET=$STRIPE_WEBHOOK_SECRET \
    STRIPE_PRICE=$STRIPE_PRICE \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
    

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download
RUN go mod verify

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN make build

COPY run.sh ./build/run.sh
COPY keystone-server-credentials.json ./build/credentials.json

FROM gcr.io/distroless/static-debian11

WORKDIR app
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/build /app/.

# Run the web service on container startup.
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json 
CMD ["/app/server"]
