FROM golang:1.13-alpine AS builder

ENV     SRC_DIR=/go/src/github.com/jmervine/noop-server
COPY    . ${SRC_DIR}
WORKDIR ${SRC_DIR}

RUN set -x; \
  go build -o ${SRC_DIR}/noop-server noop-server.go && \
  chmod 755 noop-server

FROM alpine:3
WORKDIR /src
COPY --from=builder /go/src/github.com/jmervine/noop-server/noop-server .

CMD /src/noop-server
