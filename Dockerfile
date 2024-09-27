# syntax=docker/dockerfile:1

FROM golang:1-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o operator


FROM alpine

COPY --from=build /src/operator /operator

CMD ["/operator"]