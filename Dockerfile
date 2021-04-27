FROM golang:1.16-alpine AS builder

ENV     SRC_DIR=/go/src/github.com/jmervine/noop-server
COPY    . ${SRC_DIR}
WORKDIR ${SRC_DIR}

RUN set -x; \
  go build -o ${SRC_DIR}/noop-server . && \
  chmod 755 noop-server

FROM alpine:3
WORKDIR /src
COPY --from=builder /go/src/github.com/jmervine/noop-server/noop-server .

CMD /src/noop-server
