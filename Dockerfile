# Alpine is chosen for its small footprint
FROM golang:1.17.7-alpine3.15

WORKDIR /app
COPY . ./
# Install Go Dependencies
RUN go mod download && go mod verify
# Install C compiler gcc
RUN apk add build-base
# Build Go Binary
RUN go build -o go-minitwit main.go

RUN chmod +x go-minitwit

# Install Bugsnag Panic Monitor
# It report unhandled panics to Bugsnag
RUN go install github.com/bugsnag/panic-monitor@latest

CMD ["/bin/sh","-c","panic-monitor /app/go-minitwit"]