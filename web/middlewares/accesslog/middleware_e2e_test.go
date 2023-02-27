// go:build e2e
package accesslog

import (
	"fmt"
	"gitime/web"
	"testing"
)

func TestMiddlewareBuilderE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		//fmt.Println("hello, it is me")
		ctx.Resp.Write([]byte("hello, it is me"))
	})

	server.Start(":8081")
}
