FROM golang:1.21.6 AS base
RUN apt-get install -y make

WORKDIR /opt/tt-go-json-fernet

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY ./ ./

RUN --mount=type=cache,target=/gomod-cache \
  go mod download

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    make generate

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    make build
    

FROM base AS migration
CMD ["/opt/tt-go-json-fernet/bin/migrate"]

FROM base AS server
CMD ["/opt/tt-go-json-fernet/bin/tt-go-json-fernet"]