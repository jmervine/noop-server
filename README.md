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

## mTLS Support

**Warnings:**
* You kind of need to know a bit about TLS and mTLS to use this feature.
* mTLS mode (specifically TLS) does not work on Heroku.

The mTLS support was added fairly quickly to support a test case I needed. Thus
it's had limited testing and use.

I have left my example certs in this repo under [examples/mtls](examples/mtls).

### Running the server

You need the following files:
* Self-signed server cert exported to `TLS_CERT`.
* Self-signed server key exported to `TLS_KEY`.
* Self-signed server CA chain cert exported to `MTLS_CA_CHAIN_CERT`.

> _Note: You need to export the contents of these files, not the paths. Example:_
>
> ```
> export MTLS_CA_CHAIN_CERT="`cat examples/mtls/server/ca-chain.cert.pem`"
> export TLS_CERT="`cat examples/mtls/server/localhost.cert.pem`"
> export TLS_KEY="`cat examples/mtls/server/localhost.key.pem`"
> ```

### Connecting via curl

You then need to connect using the coorosponding client cert and key and the
server ca chain cert. Example:
```
export CURL_KEY=examples/mtls/client/localhost.key.pem
export CURL_CERT=examples/mtls/client/localhost.cert.pem
export CURL_CACERT=examples/mtls/server/ca-chain.cert.pem

curl -v --cacert $CURL_CACERT --cert $CURL_CERT --key $CURL_KEY https://localhost:3000
```

### Generating the certs

I used the generate.sh script from https://github.com/nicholasjackson/mtls-go-example/.
I then copied the generated certs out for my use.
