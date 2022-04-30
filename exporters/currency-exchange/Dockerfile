# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /currency-exchange

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /currency-exchange /currency-exchange

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/currency-exchange"]