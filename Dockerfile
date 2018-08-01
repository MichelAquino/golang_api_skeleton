# Based on this image: https://hub.docker.com/_/golang/
FROM golang:latest

# Install godep
RUN go get -u github.com/golang/dep/cmd/dep

# Copy directory locally to container's directory
ADD . $GOPATH/src/github.com/michelaquino/golang_api_skeleton

# Set work directory
WORKDIR /go/src/github.com/michelaquino/golang_api_skeleton

# Install dependencies
RUN make setup

# Compile application
RUN GOOS=linux GOARCH=amd64 go build -o golang_api_skeleton main.go

# Execite application when container is started
ENTRYPOINT /go/src/github.com/michelaquino/golang_api_skeleton/golang_api_skeleton

# Expose 8080 port
EXPOSE 8080