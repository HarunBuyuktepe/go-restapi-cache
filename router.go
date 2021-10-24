package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type IRouter interface {
	InitRouter() *http.ServeMux
}

type router struct{}

func (router *router) InitRouter() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", Chain(SampleUsage, Logging()))
	serveMux.HandleFunc("/Get", Chain(GetKey, Method("GET"), Logging()))
	serveMux.HandleFunc("/Set", Chain(SetKey, Method("GET"), Logging()))
	serveMux.HandleFunc("/GetPath", Chain(GetPath, Method("GET"), Logging()))
	serveMux.HandleFunc("/SetPath", Chain(SetPath, Method("GET"), Logging()))
	serveMux.HandleFunc("/GetFrequency", Chain(GetFrequency, Method("GET"), Logging()))
	serveMux.HandleFunc("/SetFrequency", Chain(SetFrequency, Method("GET"), Logging()))
	serveMux.HandleFunc("/GetImageOfMemory", Chain(GetImageOfMemory, Method("GET"), Logging()))
	serveMux.HandleFunc("/Delete", Chain(Delete, Method("GET"), Logging()))
	serveMux.HandleFunc("/Flush", Chain(Flush, Method("GET"), Logging()))
	return serveMux
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()["key"]
	fmt.Fprintln(w, c.Get(keys[0]))
}

func SetKey(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()["key"]
	values := r.URL.Query()["value"]
	c.Set(keys[0],values[0])
	fmt.Fprintln(w, "Key value pair saved")
}

func SetFrequency(w http.ResponseWriter, r *http.Request) {
	frequency := r.URL.Query()["frequency"]
	intVar := 60
	intVar, _ = strconv.Atoi(frequency[0])
	c.SetFrequency(intVar)
	fmt.Fprintln(w, "Frequency edited",frequency[0])
}

func GetFrequency(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, c.GetFrequency())
}

func SetPath(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query()["path"]
	c.SetPath(path[0])
	fmt.Fprintln(w, "Path edited",path[0])
}

func GetPath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, c.GetPath())
}

func GetImageOfMemory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, c.GetValue())
}

func Delete(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()["key"]
	item := c.Get(keys[0])
	c.Delete(keys[0])
	fmt.Fprintln(w, keys[0],"->",item,"was deleted")
}

func Flush(w http.ResponseWriter, r *http.Request) {
	c.DeleteAll()
	fmt.Fprintln(w, "Flushed data")
}

func SampleUsage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Sample Usage\n/Get?key=XXX\n/Set?key=XXX&value=XXX")
	fmt.Fprintln(w, "/Flush\n/Delete?key=XXX")
	fmt.Fprintln(w, "/GetFrequency\n/SetFrequency?Frequency=XXX")
	fmt.Fprintln(w, "/GetPath\n/SetPath?Path=XXX")
	fmt.Fprintln(w, "/GetImageOfMemory")
}


var (
	myRouter	*router
	routerOnce	sync.Once
)

func Router() IRouter {
	routerOnce.Do(func() {
		myRouter = &router{}
	})
	return myRouter
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware check
			if r.Method != m {
				keys, _ := r.URL.Query()["key"]
				fmt.Println(keys)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
