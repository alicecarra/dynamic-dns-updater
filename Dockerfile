FROM golang:1.21 AS build-stage


WORKDIR /dynamic-dns-updater

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/dynamic-dns-updater/main.go ./

RUN go build -o /dynamic-dns-updater


FROM debian:bookworm AS build-release-stage

RUN apt-get update && apt-get upgrade
RUN apt-get install -y ca-certificates


WORKDIR /

COPY --from=build-stage /dynamic-dns-updater /dynamic-dns-updater

# CMD ["cat", "/etc/resolv.conf"]

CMD ["./dynamic-dns-updater/dynamic-dns-updater"]