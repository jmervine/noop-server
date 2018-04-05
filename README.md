# noop-server

A simple noop server that accepts everything.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)


## run

### w/ docker

```
docker run --rm -it -p 3000:3000 jmervine/noop-server:latest
```

### w/ docker compose (locally)

```
# fork and/or clone
docker-compose build
docker-compose up
```

### w/o docker

```
# fork and/or clone
go run noop-server.go
```

## use

```
jmervine@debian:~$ curl -i localhost:3000
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:06 GMT
Content-Length: 3

200 OK
jmervine@debian:~$ curl -i -X POST -d 'foo=bar' localhost:3000
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:16 GMT
Content-Length: 3

200 OK
jmervine@debian:~$ curl -i -X DELETE localhost:3000
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:24 GMT
Content-Length: 3

200 OK
jmervine@debian:~$ curl -i -H 'X-HTTP-Status: 500' -X DELETE localhost:3000
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:29 GMT
Content-Length: 22

500 Internal Server Error
jmervine@debian:~$ curl -i -H 'X-HTTP-Status: 404' -X DELETE localhost:3000/
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:41 GMT
Content-Length: 10

404 Not Found
```

## Sample using a quick start function

```
## Noop Server - Clean http listener
function listen {
  local port="$1"
  local addr="$2"

  [[ -z $port ]] && port=3000
  [[ -z $addr ]] && addr=0.0.0.0

  docker run --rm -it -p "$port:$port" -e PORT=$port -e ADDR=$addr jmervine/noop-server
}
```
