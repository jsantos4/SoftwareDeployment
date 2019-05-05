# Start the Go app build
FROM golang:latest AS build

# Copy source
WORKDIR /home/jsantos4/go/src/SoftwareDeployment
COPY . .

# Get required modules (assumes packages have been added to ./vendor)
RUN go get -d -v ./...

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /root/

# Copy files from previous build container
COPY --from=build /home/jsantos4/go/src/SoftwareDeployment ./

# Add environment variables
# ENV ...
ENV LOGGLY_TOKEN fd30ee90-c128-4f11-b5af-2b00fc3a2afc
ENV APP /home/jsantos4/go/src/SoftwareDeployment
ENV AWS_ACCESS_KEY_ID
ENV AWS_SECRET_ACCESS_KEY

# Check results
RUN env && pwd && find .

# Start the application
CMD ["./main"]
