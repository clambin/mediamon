FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ENV VERSION=$VERSION

WORKDIR /app/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-X main.version=$VERSION" \
    -o mediamon \
    cmd/mediamon/mediamon.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/mediamon /app/mediamon

RUN /usr/sbin/addgroup app
RUN /usr/sbin/adduser app -G app -D
USER app

ENTRYPOINT ["/app/mediamon"]
CMD []
