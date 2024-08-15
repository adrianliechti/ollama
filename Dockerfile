# syntax=docker/dockerfile:1

FROM golang:1 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o launcher


FROM ollama/ollama:0.3.6

RUN apt-get update && apt-get install -y \
  tini \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /src/launcher /usr/bin/launcher

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/usr/bin/launcher"]