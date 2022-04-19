FROM golang:1.18.1-alpine3.15 
WORKDIR /app
COPY go.mod go.sum ./
# Install C compiler gcc
RUN apk add build-base
# Copy files
COPY src/ src/
COPY api/ api/
# Build the binary and give executable rights
RUN go build -o go-minitwit-api api/api.go
# Tool that wait for MySQL to be ready
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x go-minitwit-api /wait

# Bugsnag panic
RUN go install github.com/bugsnag/panic-monitor@latest
CMD ["/bin/sh","-c","/wait && panic-monitor /app/go-minitwit-api"]