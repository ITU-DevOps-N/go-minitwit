FROM golang:1.17.7-alpine3.15 as base
WORKDIR /app
COPY go.mod go.sum ./
# Install Go Dependencies
RUN go mod download && go mod verify
# Install C compiler gcc
RUN apk add build-base
# Copy files
COPY src/ src/
COPY api/ api/
# Build the binary and give executable rights
RUN go build -o go-minitwit src/main.go
RUN go build -o go-minitwit-api api/api.go

RUN chmod +x go-minitwit go-minitwit-api

# Bugsnag panic
RUN go install github.com/bugsnag/panic-monitor@latest

FROM base as web
COPY --from=base /app/go-minitwit .
CMD ["/bin/sh","-c","panic-monitor /app/go-minitwit"]

FROM base as api
COPY --from=base /app/go-minitwit-api .
CMD ["/bin/sh","-c","panic-monitor /app/go-minitwit-api"]
