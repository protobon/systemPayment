FROM golang:alpine

# Add Maintainer info
LABEL maintainer="Ayrton Coelho"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup app & build folders
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /app
RUN mkdir "/build"

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

#Setup hot-reload for dev stage
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
RUN go install golang.org/x/tools/gopls@latest

ENTRYPOINT CompileDaemon -build="go build -o /build/app -buildvcs=false" -command="/build/app"
