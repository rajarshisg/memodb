# Use official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container to /memodb
WORKDIR /memodb

# Copy only the content of the /memodb directory from the host to /memodb inside the container
COPY . .

# Download Go modules (dependencies)
RUN go mod tidy

# Build the Go application
RUN go build -v -o /usr/local/bin/app .

# Specify the default command to run when the container starts
CMD ["app"]
