FROM golang:alpine

# Add Maintainer info
LABEL maintainer="Ayrton Coelho"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Download and install dependencies + build binary file
RUN go get -d -v ./...  && \
    go install -v ./...  && \
    go build -o /api

# run application
CMD [ "/api" ]