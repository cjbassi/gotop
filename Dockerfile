# Build the binary with:
# <CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .>

FROM alpine

COPY ./gotop /gotop

ENTRYPOINT ["/gotop"]
