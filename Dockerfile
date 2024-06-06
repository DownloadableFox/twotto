# Use the official Go image as the base image
FROM golang:latest

# Set environment variables
ARG COMMIT_SHA
ENV COMMIT_SHA=${COMMIT_SHA}

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies like make, google wire & go-migrate
RUN apt-get update && apt-get install -y make 
RUN go install github.com/google/wire/cmd/wire@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Generate the Go application
RUN make generate

# Change the permissions of the binary
RUN chmod +x ./bin/bot

# Set the command to run the executable
CMD ["make", "run"]