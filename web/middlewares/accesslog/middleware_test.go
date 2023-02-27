package accesslog

import (
	"fmt"
	"gitime/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it is me")
	})
	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	//req.Host = "localhost"
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
