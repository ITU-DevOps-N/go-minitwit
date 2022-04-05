FROM golang:1.17.7-alpine3.15 
WORKDIR /app
COPY go.mod go.sum ./
# Install C compiler gcc
RUN apk add build-base
# Copy files
COPY src/ src/
# Build the binary and give executable rights
RUN go build -o go-minitwit-frontend src/main.go
# Tool that wait for MySQL to be ready
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x go-minitwit-frontend /wait

# Bugsnag panic
RUN go install github.com/bugsnag/panic-monitor@latest
CMD ["/bin/sh","-c","/wait && panic-monitor /app/go-minitwit-frontend"]