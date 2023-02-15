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

# Download and install the dependencies + build binary file
RUN go get -d -v ./...  && \
    go install -v ./...  && \
    go build -o /api

EXPOSE 8081

# run application
CMD [ "/api" ]