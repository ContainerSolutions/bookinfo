############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/bookStockAPI/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/bookStockAPI
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/bookStockAPI /go/bin/bookStockAPI
WORKDIR /go/bin/
ADD static/version.txt /go/bin/static/version.txt
ADD configuration/livesettings.json /go/bin/configuration/livesettings.json
ADD swagger.yaml /go/bin/swagger.yaml
ENV BASE_URL :5550
# Run the hello binary.
ENTRYPOINT ["/go/bin/bookStockAPI"]