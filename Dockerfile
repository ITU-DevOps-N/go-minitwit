# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
FROM golang:1.17-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN go build -o /minitwit

CMD ["/minitwit"]