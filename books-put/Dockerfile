FROM golang:1.24 
# Base image: use the official Go 1.24 image from Docker Hub
# This includes the Go compiler and tools needed to build Go programs

WORKDIR /usr/src/app
# Set working directory inside the container

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download
# Download the dependencies in go.mod

COPY . .
# Copy the entire current directory into the container
RUN go build -v -o /usr/local/bin/app ./...
# Compile the Go code

CMD ["app"]
# Default command line that runs when the container starts
# Specify the name of the compiled Go binary

EXPOSE 8083
# Listen on port 3033
