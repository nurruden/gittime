package web

import (
	"net"
	"net/http"
)

var _ Server = &HTTPServer{}

// type HandleFunc func(ctx *Context)
type HandleFunc func(ctx *Context)
type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(method string, path string, handleFunc HandleFunc)
}
type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	router
	mdls []Middleware
}

// type HTTPSServer struct {
//
// }
//
//	func (h *HTTPServer) addRoute(method string,path string,handleFunc HandleFunc){
//		//panic("implement me")
//
// }
func NewHTTPServerV1(mdls ...Middleware) *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
		mdls:   mdls,
	}

}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	root := h.serve
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}
	root(ctx)
}
func (h *HTTPServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		_, _ = ctx.Resp.Write([]byte("Not found"))
		return
	}
	ctx.PathParams = info.pathParams
	ctx.MatchedRoute = info.n.route
	info.n.handler(ctx)
}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}
func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServer) Options(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
	//panic("implement me")
}
