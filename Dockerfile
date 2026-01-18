# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /src

# Copy sources and download modules
COPY go.mod go.sum .
RUN go mod download
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o /bin/helloserver

# Final stage
FROM scratch

# Copy the static binary
COPY --from=build /bin/helloserver /

# The service listens on port 8080. 
# The user can override this with the -addr flag.
EXPOSE 8080

ENTRYPOINT ["/helloserver"]
