FROM golang:1.13-alpine

ENV     SRC_DIR=/go/src/github.com/jmervine/noop-server
COPY    . ${SRC_DIR}
WORKDIR ${SRC_DIR}

RUN set -x; \
  go build -o ${SRC_DIR}/noop-server noop-server.go && \
  chmod 755 noop-server

CMD ${SRC_DIR}/noop-server
