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
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})
	h.Post("/form", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello,%s", ctx.Req.URL.Path)))
	})
	h.Get("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("wrong id input"))
			return
		}

		ctx.Resp.Write([]byte(fmt.Sprintf("hello,%d", id)))
	})
	type User struct {
		Name string `json:"name"`
	}
	h.Get("/user/123", func(ctx *Context) {
		ctx.RespJSON(202, User{
			Name: "tom",
		})
	})
	h.Start(":8081")
}
