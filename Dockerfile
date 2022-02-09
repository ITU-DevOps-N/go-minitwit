# syntax=docker/dockerfile:1

# link https://docs.docker.com/language/golang/build-images/

# We use golang version 1.17 running a lightweight linux distribution (alpine)
FROM golang:1.17-alpine

# Working directory inside the container, and copy over the static files
WORKDIR /app

# Go convention, files that track module(s) which includes our packages
COPY go.mod ./
COPY go.sum ./

# Install the Go modules inside the container
RUN go mod download

# Add source code to the container
COPY . ./



ENV PATH="/scripts:${PATH}:/sbin"

ADD /itu-minitwit /app
WORKDIR /app

RUN apk update
RUN apk add --update --no-cache --virtual .tmp gcc libc-dev linux-headers
RUN apk add python3
RUN apk add py3-pip
RUN apk add sqlite sqlite-dev 

RUN pip3 install -r requirements.txt 

CMD python3 minitwit.py