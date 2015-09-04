FROM     golang:1.5
ENV      SRC_DIR=/go/src/github.com/jmervine/noop-server
ADD      . ${SRC_DIR}
WORKDIR  ${SRC_DIR}
RUN      go build -o ${SRC_DIR}/noop-server noop-server.go && chmod 755 noop-server
CMD      ${SRC_DIR}/noop-server
