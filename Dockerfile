############################
# STEP 1 build application
############################
FROM golang:1.16.6-alpine3.13 AS builder

# Get certificates
RUN apk update \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY config ./config
COPY utils.go .
COPY main.go .
RUN go build -o main .

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main /

ENTRYPOINT ["/main"]
