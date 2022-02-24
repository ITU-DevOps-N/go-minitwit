# Alpine is chosen for its small footprint
FROM golang:alpine3.15

# Heavy build alternative:

# WORKDIR /app
# COPY . ./
# Install Go Dependencies
# RUN go mod download && go mod verify
# # Install C compiler gcc
# RUN apk add build-base
# Build Go Binary
# RUN go build -o minitwit main.go

# Light build alternative:
# Requires minitwit binary in root directory
COPY minitwit ./

RUN chmod +x minitwit

CMD ["./minitwit"]