package web

import (
	"fmt"
	"net/http"
	"testing"
)

//func TestServer(t *testing.T)  {
//	var h Server
//	http.ListenAndServe("8081",h)
//
//}

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("first before")
				next(ctx)
				fmt.Println("first after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("second before")
				next(ctx)
				fmt.Println("second after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("third  interrupt")

			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("cannot see")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
