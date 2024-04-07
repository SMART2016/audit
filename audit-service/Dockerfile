# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM debian:buster as builder
RUN apt-get update && apt-get install -y wget

ENV GOLANG_VERSION 1.22.2
RUN wget https://go.dev/dl/go1.22.2.linux-arm64.tar.gz \
    && tar -C /usr/local -xzf go1.22.2.linux-arm64.tar.gz \
    && rm go1.22.2.linux-arm64.tar.gz
ENV PATH $PATH:/usr/local/go/bin

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY .. ./

# Build the binary.
# -o myapp specifies the output filename (myapp)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o audit-service
#RUN go build -v -o audit-service

# Use the official Debian slim image for a lean production container.
FROM debian:buster-slim
WORKDIR /root/

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/audit-service .

# Service listens on port 8080.
EXPOSE 8080

# Run the web service on container startup.
CMD ["./audit-service"]
