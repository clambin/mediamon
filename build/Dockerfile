FROM golang:1.15 AS builder

WORKDIR /build

COPY . ./

ARG BUILD_VERSION=development
RUN CGO_ENABLED=0 go build -ldflags "-X covid19/internal/version.BuildVersion=$BUILD_VERSION" cmd/mediamon/mediamon.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/mediamon /app

EXPOSE 8080
ENTRYPOINT ["/app/mediamon"]
CMD []
