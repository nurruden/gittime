package web

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
}

type node struct {
	route      string
	path       string
	children   map[string]*node
	paramChild *node
	starChild  *node
	handler    HandleFunc
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

// func (n *node)childOf(seg string) (*node, bool){
//
// }
func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	//panic("implement me")
	if path == "" {
		panic("web: path cannot be null string")
	}
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
			//handler: handleFunc,
		}
		r.trees[method] = root
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: path cannot be ended with /")
	}
	if path[0] != '/' {
		panic("web: path must be started with /")
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: route conflict, duplicate path")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}
	//path = path[1:]
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web: path cannot contain continuously / ")
		}
		children := root.childOrCreate(seg)
		root = children
	}
	if root.handler != nil {
		//panic("web: route conflict, duplicate path")
		panic(fmt.Sprintf("web: route conflict, duplicate path [%s]", path))
	}
	root.handler = handleFunc
	root.route = path

}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
	//return root, root.handler != nil
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}

func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: Not allowed register path param and wildcard at the same time, path param exist")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}
	if seg == "*" {
		if n.paramChild != nil {
			panic("web: Not allowed register path param and wildcard at the same time, wildcard exist")
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}
	if n.children == nil {
		n.children = map[string]*node{}

	}

	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}
