#!/usr/bin/env bash

function do_curl {
  curl -H 'X-NoopServerFlags:echo;sleep=500ms;status=301;host=http://test.host' -d 'foo=bar' localhost:3000/load/$1
  echo
}

while true
do
  do_curl "1"
  do_curl "2"
  do_curl "3"
  sleep .05
done