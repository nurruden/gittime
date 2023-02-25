package web

import (
	"net"
	"net/http"
)
var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(method string,path string, handleFunc HandleFunc)
}

type HTTPServer struct {
	router
}

//type HTTPSServer struct {
//
//}
//func (h *HTTPServer) addRoute(method string,path string,handleFunc HandleFunc){
//	//panic("implement me")
//
//}

func NewHTTPServer() *HTTPServer{
	return &HTTPServer{
		router: newRouter(),
	}
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req: request,
		Resp: writer,
	}
	h.serve(ctx)
}
func (h *HTTPServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method,ctx.Req.URL.Path)
	if !ok || info.n.handler == nil{
		ctx.Resp.WriteHeader(404)
		_,_ = ctx.Resp.Write([]byte("Not found"))
		return
	}
	ctx.PathParams = info.pathParams
	info.n.handler(ctx)
}

func (h *HTTPServer) Get(path string,handleFunc HandleFunc){
	h.addRoute(http.MethodGet,path,handleFunc)
}
func (h *HTTPServer) Post(path string,handleFunc HandleFunc){
	h.addRoute(http.MethodPost,path,handleFunc)
}

func (h *HTTPServer) Options(path string,handleFunc HandleFunc){
	h.addRoute(http.MethodOptions,path,handleFunc)
}

func (h *HTTPServer) Start(addr string) error {
	l,err :=net.Listen("tcp",addr)
	if err != nil{
		return err
	}
	return http.Serve(l,h)
	//panic("implement me")
}
