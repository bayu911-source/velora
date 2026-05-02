# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 go build -o /usr/local/bin/velora ./...

FROM gcr.io/distroless/base-debian11
COPY --from=build /usr/local/bin/velora /usr/local/bin/velora
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/velora", "server"]
