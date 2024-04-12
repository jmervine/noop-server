FROM golang:1-alpine AS builder

ENV     SRC_DIR=/go/src/github.com/jmervine/noop-server
COPY    . ${SRC_DIR}
WORKDIR ${SRC_DIR}

RUN set -x; \
  mkdir -p bin && \
  go build -o ${SRC_DIR}/bin/noop-server ./cmd/noop-server/ && \
  chmod 755 ${SRC_DIR}/bin/noop-server

FROM alpine:3
WORKDIR /src
COPY --from=builder /go/src/github.com/jmervine/noop-server/bin/noop-server .

CMD /src/noop-server
