# env

[![GoDoc](https://godoc.org/gopkg.in/jmervine/env.v1?status.png)](https://godoc.org/gopkg.in/jmervine/env.v1) [![Build Status](https://travis-ci.org/jmervine/env.svg?branch=master)](https://travis-ci.org/jmervine/env)


Simple configuration utility around os.{Get,Set}env

## usage

```

PACKAGE DOCUMENTATION

package env
    import "github.com/jmervine/env"

    env is a simple package for loading configuration, based loosly on
    Ruby's `dotenv` gem.

    Example:

	package main

	import (
		"github.com/jmervine/env"

		"fmt"
	)

	func main() {
		var err error
		err = env.Load("example.env")
		if err != nil {
	  	panic(err)
		}

		env.PanicOnRequire = true

		d, _ := env.Require("DATABASE_URL")
		var (
			dburl   = d
			ignored = env.GetOrSetBool("IGNORED", true)
			debug   = env.GetBool("DEBUG")
			addr    = env.GetString("ADDR")
			port    = env.GetOrSetInt("PORT", 3000)
		)

		fmt.Printf("dburl   ::: %s\n", dburl)
		fmt.Printf("ignored ::: %v\n", ignored)
		fmt.Printf("debug   ::: %v\n", debug)
		fmt.Printf("addr    ::: %s\n", addr)
		fmt.Printf("port    ::: %d\n", port)
	}

VARIABLES

var PanicOnRequire = false
    PanicOnRequire forces panics when Require- methods fail

FUNCTIONS

func Get(key string) string
    Get gets a key and returns a string

func GetBool(key string) bool
    GetBool gets a key and sets to true, false or nil using the Truthy and
    Falsey variables

func GetDuration(key string) time.Duration

func GetFloat32(key string) float32
    GetFloat32 gets a key and returns an float32

func GetFloat64(key string) float64
    GetFloat64 gets a key and returns an float64

func GetInt(key string) int
    GetInt gets a key and returns an int

func GetInt32(key string) int32
    GetInt32 gets a key and returns an int32

func GetInt64(key string) int64
    GetInt64 gets a key and returns an int64

func GetOrSet(key string, val interface{}) string
    GetOrSet gets a key and returns a string or set's the default

func GetOrSetBool(key string, val bool) bool

func GetOrSetDuration(key string, val time.Duration) time.Duration

func GetOrSetFloat32(key string, val float32) float32

func GetOrSetFloat64(key string, val float64) float64

func GetOrSetInt(key string, val int) int

func GetOrSetInt32(key string, val int32) int32

func GetOrSetInt64(key string, val int64) int64

func GetString(key string) string
    GetString is an alias to Get

func Load(filenames ...string) error
    Load loads a file containing standard os environment key/value pairs,
    doesn't override currently set variables

    e.g.: .env

	PORT=3000
	ADDR=0.0.0.0
	DEBUG=true

func Overload(filenames ...string) error
    Overload does the same thing as Load, but overrides existing variables

func Require(key string) (val string, err error)
    Require gets a key and returns a string or an error if it's set to "" in
    os.Getenv

func RequireBool(key string) (bool, error)

func RequireDuration(key string) (time.Duration, error)

func RequireFloat32(key string) (float32, error)

func RequireFloat64(key string) (float64, error)

func RequireInt(key string) (int, error)

func RequireInt32(key string) (int32, error)

func RequireInt64(key string) (int64, error)

func RequireString(key string) (string, error)

func Set(key string, val interface{})
    Set sets via an interface

func SetMap(m map[string]interface{})
    SetMap iterates over a map and sets keys to values

```

## LICENSE

```
Copyright (c) 2015 Joshua Mervine

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
