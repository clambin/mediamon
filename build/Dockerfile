FROM golang:1.15 AS builder

WORKDIR /build

COPY . ./

ARG BUILD_VERSION=development
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/clambin/mediamon/internal/version.BuildVersion=$BUILD_VERSION" cmd/github.com/clambin/mediamon/github.com/clambin/mediamon.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/github.com/clambin/mediamon /app

EXPOSE 8080
ENTRYPOINT ["/app/github.com/clambin/mediamon"]
CMD []
