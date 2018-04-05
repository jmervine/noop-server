// env is a simple package for loading configuration, based loosly on
// Ruby's `dotenv` gem.
//
// Example:
//
//    package main
//
//    import (
//    	"github.com/jmervine/env"
//
//    	"fmt"
//    )
//
//    func main() {
//    	var err error
//    	err = env.Load("example.env")
//    	if err != nil {
//      	panic(err)
//    	}
//
//    	env.PanicOnRequire = true
//
//    	d, _ := env.Require("DATABASE_URL")
//    	var (
//    		dburl   = d
//    		ignored = env.GetOrSetBool("IGNORED", true)
//    		debug   = env.GetBool("DEBUG")
//    		addr    = env.GetString("ADDR")
//    		port    = env.GetOrSetInt("PORT", 3000)
//    	)
//
//    	fmt.Printf("dburl   ::: %s\n", dburl)
//    	fmt.Printf("ignored ::: %v\n", ignored)
//    	fmt.Printf("debug   ::: %v\n", debug)
//    	fmt.Printf("addr    ::: %s\n", addr)
//    	fmt.Printf("port    ::: %d\n", port)
//    }
package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	dotenv "github.com/joho/godotenv"
)

// PanicOnRequire forces panics when Require- methods fail
var PanicOnRequire = false

// Load loads a file containing standard os environment key/value pairs,
// doesn't override currently set variables
//
// e.g.: .env
//
//     PORT=3000
//     ADDR=0.0.0.0
//     DEBUG=true
//
func Load(filenames ...string) error {
	return dotenv.Load(filenames...)
}

// Overload does the same thing as Load, but overrides existing variables
func Overload(filenames ...string) error {
	env, err := dotenv.Read(filenames...)
	if err != nil {
		return err
	}

	for key, val := range env {
		os.Setenv(key, val)
	}

	return nil
}

// Set sets via an interface
func Set(key string, val interface{}) {
	os.Setenv(key, toString(val))
}

// SetMap iterates over a map and sets keys to values
func SetMap(m map[string]interface{}) {
	for key, val := range m {
		Set(key, val)
	}
}

// Get gets a key and returns a string
func Get(key string) string {
	return os.Getenv(key)
}

// GetString is an alias to Get
func GetString(key string) string {
	return Get(key)
}

// Require gets a key and returns a string or an error if it's set to "" in
// os.Getenv
func Require(key string) (val string, err error) {
	val = Get(key)
	if val == "" {
		err = onError(fmt.Errorf("missing required string from %s", key))
	}

	return val, err
}

func RequireString(key string) (string, error) {
	return Require(key)
}

// GetOrSet gets a key and returns a string or set's the default
func GetOrSet(key string, val interface{}) string {
	str := Get(key)

	if str != "" {
		return str
	}

	v := toString(val)
	Set(key, v)

	return v
}

func RequireDuration(key string) (time.Duration, error) {
	str, err := Require(key)
	if err != nil {
		d := new(time.Duration)
		return *d, onError(fmt.Errorf("missing required duration from %s", key))
	}

	return toDur(str), nil
}

func GetDuration(key string) time.Duration {
	return toDur(Get(key))
}

func GetOrSetDuration(key string, val time.Duration) time.Duration {
	str := Get(key)
	if str != "" {
		return toDur(str)
	}

	Set(key, val)
	return val
}

// GetInt gets a key and returns an int
func GetInt(key string) int {
	return toInt(Get(key))
}

func GetOrSetInt(key string, val int) int {
	str := Get(key)
	if str != "" {
		return toInt(str)
	}
	Set(key, val)
	return val
}

func RequireInt(key string) (int, error) {
	str := Get(key)
	if str != "" {
		return toInt(str), nil
	}

	return int(0), onError(fmt.Errorf("missing required int from %s", key))
}

// GetInt32 gets a key and returns an int32
func GetInt32(key string) int32 {
	return toInt32(Get(key))
}

func GetOrSetInt32(key string, val int32) int32 {
	str := Get(key)
	if str != "" {
		return toInt32(str)
	}
	Set(key, val)
	return val
}

func RequireInt32(key string) (int32, error) {
	str := Get(key)
	if str == "" {
		return int32(0), onError(fmt.Errorf("missing required int32 from %s", key))
	}
	return toInt32(str), nil
}

// GetInt64 gets a key and returns an int64
func GetInt64(key string) int64 {
	return toInt64(Get(key))
}

func GetOrSetInt64(key string, val int64) int64 {
	str := Get(key)
	if str != "" {
		return toInt64(str)
	}
	Set(key, val)
	return val
}

func RequireInt64(key string) (int64, error) {
	str := Get(key)
	if str == "" {
		return int64(0), onError(fmt.Errorf("missing required int64 from %s", key))
	}
	return toInt64(str), nil
}

// GetFloat32 gets a key and returns an float32
func GetFloat32(key string) float32 {
	return toFloat32(Get(key))
}

func GetOrSetFloat32(key string, val float32) float32 {
	str := Get(key)
	if str != "" {
		return toFloat32(str)
	}
	Set(key, val)
	return val
}

func RequireFloat32(key string) (float32, error) {
	str := Get(key)
	if str == "" {
		return float32(0), onError(fmt.Errorf("missing required float32 from %s", key))
	}
	return toFloat32(str), nil
}

// GetFloat64 gets a key and returns an float64
func GetFloat64(key string) float64 {
	return toFloat64(Get(key))
}

func GetOrSetFloat64(key string, val float64) float64 {
	str := Get(key)
	if str != "" {
		return toFloat64(str)
	}
	Set(key, val)
	return val
}

func RequireFloat64(key string) (float64, error) {
	str := Get(key)
	if str == "" {
		return float64(0), onError(fmt.Errorf("missing required float64 from %s", key))
	}
	return toFloat64(str), nil
}

// GetBool gets a key and sets to true, false or nil using the Truthy and Falsey
// variables
func GetBool(key string) bool {
	return toBool(Get(key))
}

func GetOrSetBool(key string, val bool) bool {
	str := Get(key)
	if str != "" {
		return toBool(str)
	}
	Set(key, val)
	return val
}

func RequireBool(key string) (bool, error) {
	str := Get(key)
	if str == "" {
		return false, onError(fmt.Errorf("missing required bool from %s", key))
	}
	return toBool(str), nil
}

// HELPERS
func toBool(val string) bool {
	b, _ := strconv.ParseBool(val)
	return b
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func toDur(str string) time.Duration {
	dur, _ := time.ParseDuration(str)
	return dur
}

func toInt(val string) int {
	i, _ := strconv.ParseInt(val, 10, 16)
	return int(i)
}

func toInt32(val string) int32 {
	i, _ := strconv.ParseInt(val, 10, 32)
	return int32(i)
}

func toInt64(val string) int64 {
	i, _ := strconv.ParseInt(val, 10, 64)
	return int64(i)
}

func toFloat32(val string) float32 {
	i, _ := strconv.ParseFloat(val, 32)
	return float32(i)
}

func toFloat64(val string) float64 {
	i, _ := strconv.ParseFloat(val, 64)
	return float64(i)
}

func onError(e error) error {
	if PanicOnRequire {
		panic(e)
	}
	return e
}
