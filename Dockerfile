########################### Build Base ###########################
FROM golang:1.16-alpine as build_base

RUN apk add git make

# Add the current directory to be build
WORKDIR /app

######## Environment variables ########
# Force the go compiler to use modules
ENV GO111MODULE=on

# Disable Go proxy
ENV GOPROXY=direct

######## Install dependencies ########
COPY go.mod .
COPY go.sum .
COPY Makefile .
RUN make setup

############################# Builder ############################
FROM build_base AS builder

COPY . /app

# Build the image for Linux instances
RUN make build

############################# Runner #############################
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

COPY --from=builder /app/golang_api_skeleton /app/golang_api_skeleton

ENTRYPOINT ["/app/golang_api_skeleton", "api"]
