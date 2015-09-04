# noop-server

A simple noop server that accepts everything.

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

OK
jmervine@debian:~$ curl -i -X POST -d 'foo=bar' localhost:3000
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:16 GMT
Content-Length: 3

OK
jmervine@debian:~$ curl -i -X DELETE localhost:3000
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:24 GMT
Content-Length: 3

OK
jmervine@debian:~$ curl -i -X DELETE localhost:3000/500
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:29 GMT
Content-Length: 22

Internal Server Error
jmervine@debian:~$ curl -i -X DELETE localhost:3000/404
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Sep 2015 16:56:41 GMT
Content-Length: 10

Not Found
```

