

package web

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	//s := NewHTTPServer()
	//s.Get("/", func(ctx *Context) {
	//	ctx.Resp.Write([]byte("hello, world"))
	//})
	//s.Get("/user", func(ctx *Context) {
	//	ctx.Resp.Write([]byte("hello, user"))
	//})
	//
	//s.Start(":8081")
	h := NewHTTPServer()
	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})
	h.Get("/order/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s",ctx.Req.URL.Path)))
	})
	h.Start(":8081")
}