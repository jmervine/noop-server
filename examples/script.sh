#!/usr/bin/env bash
echo "$ curl http://localhost:3000/"
curl http://localhost:3000/

echo "$ curl http://localhost:3000/endpoint/1"
curl http://localhost:3000/endpoint/1

echo "$ curl http://localhost:3000/endpoint/2"
curl http://localhost:3000/endpoint/2

echo "$ curl http://localhost:3000/endpoint/3"
curl http://localhost:3000/endpoint/3

echo "$ curl http://localhost:3000/endpoint/WILD"
curl http://localhost:3000/endpoint/WILD

echo "$ curl http://localhost:3000/"
curl http://localhost:3000/