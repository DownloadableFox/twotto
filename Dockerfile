# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Download dependencies like make, google wire & go-migrate
RUN apt-get update && apt-get install -y make 
RUN go get github.com/google/wire/cmd/wire
RUN go get -u -d github.com/golang-migrate/migrate/cmd/migrate

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN make build

# Expose the port that the application listens on
EXPOSE 8080

# Set the command to run the executable
CMD ["make", "bootstrap"]