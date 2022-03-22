# Alpine is chosen for its small footprint
FROM golang:1.17.7-alpine3.15

WORKDIR /app
COPY . ./
# Install Go Dependencies
RUN go mod download && go mod verify
# Install C compiler gcc
RUN apk add build-base
# Build Go Binary
RUN go build -o go-minitwit src/main.go

RUN chmod +x go-minitwit

# Install Bugsnag Panic Monitor
# It report unhandled panics to Bugsnag
RUN go install github.com/bugsnag/panic-monitor@latest

# Tool that wait for MySQL to be ready
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

CMD ["/bin/sh","-c","/wait && panic-monitor /app/go-minitwit"]
