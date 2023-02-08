package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	testRoutes := []struct{
		method string
		path string
	}{
		{
			method: http.MethodGet,
			path: "/user/home",
		},
		{
			method: http.MethodDelete,
			path: "/",
		},
		{
			method: http.MethodGet,
			path: "/",
		},
		{
			method: http.MethodGet,
			path: "/user",
		},
		{
			method: http.MethodGet,
			path: "/order/detail",
		},
		{
			method: http.MethodGet,
			path: "/order/*",
		},
		{
			method: http.MethodGet,
			path: "/order/detail/:id",
		},
		//{
		//	method: http.MethodGet,
		//	path: "/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path: "/*/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path: "/*/abc",
		//},
		//{
		//	method: http.MethodGet,
		//	path: "/*/abc/*",
		//},
		{
			method: http.MethodPost,
			path: "/order/create",
		},
		{
			method: http.MethodPost,
			path: "/login",
		},
	}
	mockHandler := func(ctx *Context){}
	r:=newRouter()
	for _,route := range testRoutes{
		r.addRoute(route.method,route.path,mockHandler)
	}
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path: "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
													},
					"order": &node{
						path: "order",
						//handler: mockHandler,
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path: ":id",
									handler: mockHandler,
								},
							},

						},
						starChild: &node{
							path: "*",
							handler: mockHandler,

						},
					},
							},
						},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						//handler: mockHandler,
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login":{
						path: "login",
						handler: mockHandler,

					},
				},

						},
		},
	}
	msg,ok := wantRouter.equal(&r)
	assert.True(t,ok,msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"",mockHandler)

	},"web: path must be started with /")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"/a/b/c/",mockHandler)

	},"web: path cannot be ended with /")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"/a/b//c/",mockHandler)

	},"web: path cannot contain continuously / ")
	r = newRouter()
	r.addRoute(http.MethodGet,"/",mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"/",mockHandler)
	},"web: route conflict, duplicate path")
	r = newRouter()
	r.addRoute(http.MethodGet,"/a/b/c",mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"/a/b/c",mockHandler)
	},"web: route conflict, duplicate path /a/b/c")
	r = newRouter()
	r.addRoute(http.MethodGet,"/a/*",mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"/a/:id",mockHandler)
	},"web: Not allowed register path param and wildcard at the same time, path param exist")

	r = newRouter()
	r.addRoute(http.MethodGet,"/a/:id",mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet,"/a/*",mockHandler)
	},"web: Not allowed register path param and wildcard at the same time, wildcard exist")

}

func (n *node)equal(y *node)(string,bool){
	//if y == nil{
	//	return "destination is nil", false
	//}
	if n.path != y.path{
		return fmt.Sprint("node mismatch"),false
	}
	if len(n.children) != len(y.children){
		return fmt.Sprint("children length mismatch"),false
	}
	if n.starChild != nil{
		msg,ok := n.starChild.equal(y.starChild)
		if !ok{
			return msg,ok
		}
	}
	if n.paramChild != nil{
		msg,ok := n.paramChild.equal(y.paramChild)
		if !ok{
			return msg,ok
		}
	}
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler{
		return fmt.Sprintf("handler unequale"),false
	}
	for path, c := range n.children{
		dst,ok := y.children[path]
		if !ok {
			return fmt.Sprintf("children %s not exist",path),false
		}
		msg,ok := c.equal(dst)
		if !ok{
			return msg,false
		}
	}
	return "",true
}

func (r *router)equal(y *router)(string,bool){
	for k,v := range r.trees{
		dst,ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("cannot find http method"),false
		}
		msg,equal := v.equal(dst)
		if !equal{
			return msg,false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T){
	testRoute := []struct{
		method string
		path string
	}{
		{
			method: http.MethodGet,
			path: "/user/home",
		},
		{
			method: http.MethodGet,
			path: "/",
		},
		{
			method: http.MethodGet,
			path: "/user",
		},
		{
			method: http.MethodGet,
			path: "/order/detail",
		},
		{
			method: http.MethodGet,
			path: "/order/*",
		},
		{
			method: http.MethodPost,
			path: "/order/create",
		},
		{
			method: http.MethodPost,
			path: "/login",
		},
		{
			method: http.MethodPost,
			path: "/login/:username",
		},
	}
	r := newRouter()
	var mockHandler HandleFunc= func(ctx *Context) {
		
	}
	
	for _,route := range testRoute{
		r.addRoute(route.method, route.path,mockHandler)
	}
	testCases:=[]struct{
		name string
		method string
		path string
		wantFound bool
		info *matchInfo
	}{
		{
			name: "method not found",
			method: http.MethodOptions,
			path: "/order/detail",
			wantFound: false,

		},
		{
			name: "order detail",
			method: http.MethodGet,
			path: "/order/detail",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					handler: mockHandler,
					path: "detail",
				},
			},

		},
		{
			name: "order start",
			method: http.MethodGet,
			path: "/order/abc",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					handler: mockHandler,
					path: "*",
				},
			},

		},
		{
			name: "order",
			method: http.MethodGet,
			path: "/order",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					//handler: mockHandler,
					path: "order",
					children: map[string]*node{
						"detail": &node{
							handler: mockHandler,
							path: "detail",
						},
					},
				},
			},

		},
		{
			name: "root",
			method: http.MethodDelete,
			path: "/",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					handler: mockHandler,
					path: "/",
					children: map[string]*node{
						"order": &node{
							path: "order",
							children: map[string]*node{
								"detail":&node{
									handler: mockHandler,
									path: "detail",
								},
							},
						},
					},
				},
			},

		},
		{
			name: "Path not found",
			method: http.MethodDelete,
			path: "/aaa",
		},
		{
			name: "root",
			method: http.MethodDelete,
			path: "/",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					path:       "/",
					handler:    mockHandler,
				},
			},
		},
		{
			name: "login username",
			method: http.MethodPost,
			path: "/login/daming",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					path: ":username",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"username":"daming",
				},
			},
		},
	}
	for _,tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			info,found := r.findRoute(tc.method,tc.path)
			assert.Equal(t, tc.wantFound,found)
			if !found{
				return
			}
			assert.Equal(t, tc.info.pathParams,info.pathParams)
			msg,ok := tc.info.n.equal(info.n)
			assert.True(t, ok,msg)
			//assert.Equal(t, tc.wantNode.path,n.path)
			//assert.Equal(t, tc.wantNode.children,n.children)
			//nHandler := reflect.ValueOf(n.handler)
			//yHandler := reflect.ValueOf(tc.wantNode.handler)
			//assert.True(t, nHandler == yHandler)
		})
	}
}