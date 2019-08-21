package gmf

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"path"
	"unsafe"

	//	"github.com/julienschmidt/httprouter"
)
const (
	// AbortIndex .
	AbortIndex = math.MaxInt8 / 2
)
//engine的结构体形式，有个基本的路由和路由组
type Engine struct {
	*RouterGroup
	router *Router

}
type Router struct {
	maps map[string]*node
	*RouterGroup
}

type node struct {
	path      string
	wildChild bool
	handle    Handle

}
//启动http服务
//func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	engine.router.ServeHTTP(w, req)
//}
//TO DO
func (r *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	    path := req.URL.Path

		if len(path) > 1 && path[len(path)-1] == '/' {
			req.URL.Path = path[:len(path)-1]
		} else {
			req.URL.Path = path + "/"
		}
		http.Redirect(w, req, req.URL.String(), 200)
		return
	}
//其实就是listenandserve的简化
func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
// New一个Engine的方法
func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{nil, "", nil, engine}
	engine.router = NewRouter()
	return engine
}
func NewRouter() *Router {
	return &Router{
	}
}
// HandlerFunc的方法定义，怎么用主要看Context .
type HandlerFunc func(*Context)
// H .
type H map[string]interface{}
// ErrorMsg的结构，json格式展示 .
type ErrorMsg struct {
Message string      `json:"msg"`
Meta    interface{} `json:"meta"`
}
//核心 Context基本按gin的样式 .
type Context struct {
Req      *http.Request
Writer   http.ResponseWriter
Keys     map[string]interface{}//和路由无关
Errors   []ErrorMsg
Params   Params
ChildParams Param
handlers []HandlerFunc
engine   *Engine
index    int8
}

// RouterGroup .注意HandlerFunc是个数组。prefix时前缀
type RouterGroup struct {
Handlers []HandlerFunc
prefix   string
parent   *RouterGroup
engine   *Engine
}
// Group .HandlerFunc用...来确保都能用
func (group *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := path.Join(group.prefix, component)
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		parent:   group,
		prefix:   prefix,
		engine:   group.engine,
	}
}
func (group *RouterGroup) createContext(w http.ResponseWriter, req *http.Request, params Params, handlers []HandlerFunc) *Context {
	return &Context{
		Writer:   w,
		Req:      req,
		index:    -1,
		engine:   group.engine,
		Params:   params,
		handlers: handlers,
	}
}
// 对Get简化操作
func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	group.Handle("GET", path, handlers)
}


func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	group.Handle("POST", path, handlers)
}

func (group *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	group.Handle("DELETE", path, handlers)
}


func (group *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	group.Handle("PATCH", path, handlers)
}


func (group *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	group.Handle("PUT", path, handlers)
}

func (group *RouterGroup) WEBSOCKET(path string, config WebSocketConfig) {
	group.WsHandle("WEBSOCKET", path, config)
}

// Handle .
func (group *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	p = path.Join(group.prefix, p)
	handlers = group.combineHandlers(handlers)

      group.engine.router.RouteHandle(method ,p, func(w http.ResponseWriter, req *http.Request, params Params) {
		group.createContext(w, req, params, handlers).Next()
})

}

//TO DO
func (r *Router) RouteHandle(method, path string, handle Handle) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.maps == nil {
		r.maps = make(map[string]*node)
	}

	root := r.maps[method]

	root.addRoute(path, handle)
}
//TO DO
func (n *node) addRoute(path string, handle Handle) {
	fullPath := path
	fmt.Println(fullPath)
	//
	//
	//	for {
	//
	//
	//		// Split edge
	//		if i < len(n.path) {
	//			child := node{
	//				path:      n.path[i:],
	//				wildChild: n.wildChild,
	//				nType:     static,
	//				indices:   n.indices,
	//				children:  n.children,
	//				handle:    n.handle,
	//				priority:  n.priority - 1,
	//			}
	//
	//			// Update maxParams (max of all children)
	//			for i := range child.children {
	//				if child.children[i].maxParams > child.maxParams {
	//					child.maxParams = child.children[i].maxParams
	//				}
	//			}
	//
	//			n.children = []*node{&child}
	//			// []byte for proper unicode char conversion, see #65
	//			n.indices = string([]byte{n.path[i]})
	//			n.path = path[:i]
	//			n.handle = nil
	//			n.wildChild = false
	//		}
	//
	//		// Make new node a child of this node
	//		if i < len(path) {
	//			path = path[i:]
	//
	//			if n.wildChild {
	//				n = n.children[0]
	//				n.priority++
	//
	//				// Update maxParams of the child node
	//				if numParams > n.maxParams {
	//					n.maxParams = numParams
	//				}
	//				numParams--
	//
	//				// Check if the wildcard matches
	//				if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
	//				// Check for longer wildcard, e.g. :name and :names
	//					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
	//					continue walk
	//				} else {
	//					// Wildcard conflict
	//					var pathSeg string
	//					if n.nType == catchAll {
	//						pathSeg = path
	//					} else {
	//						pathSeg = strings.SplitN(path, "/", 2)[0]
	//					}
	//					prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
	//					panic("'" + pathSeg +
	//						"' in new path '" + fullPath +
	//						"' conflicts with existing wildcard '" + n.path +
	//						"' in existing prefix '" + prefix +
	//						"'")
	//				}
	//			}
	//
	//			c := path[0]
	//
	//			// slash after param
	//			if n.nType == param && c == '/' && len(n.children) == 1 {
	//				n = n.children[0]
	//				n.priority++
	//				continue walk
	//			}
	//
	//			// Check if a child with the next path byte exists
	//			for i := 0; i < len(n.indices); i++ {
	//				if c == n.indices[i] {
	//					i = n.incrementChildPrio(i)
	//					n = n.children[i]
	//					continue walk
	//				}
	//			}
	//
	//			// Otherwise insert it
	//			if c != ':' && c != '*' {
	//				// []byte for proper unicode char conversion, see #65
	//				n.indices += string([]byte{c})
	//				child := &node{
	//					maxParams: numParams,
	//				}
	//				n.children = append(n.children, child)
	//				n.incrementChildPrio(len(n.indices) - 1)
	//				n = child
	//			}
	//			n.insertChild(numParams, path, fullPath, handle)
	//			return
	//
	//		} else if i == len(path) { // Make node a (in-path) leaf
	//			if n.handle != nil {
	//				panic("a handle is already registered for path '" + fullPath + "'")
	//			}
	//			n.handle = handle
	//		}
	//		return
	//	}
	//} else { // Empty tree
	//	n.insertChild(numParams, path, fullPath, handle)
	//	n.nType = root
	//}
}


// 给Context里面的key设置变量内容
func (c *Context) Set(key string, item interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = item
}
// 读取Context里key的变量内容
// 没有就panic
func (c *Context) Get(key string) interface{} {
	var ok bool
	var item interface{}
	if c.Keys != nil {
		item, ok = c.Keys[key]
	} else {
		item, ok = nil, false
	}
	if !ok || item == nil {
		log.Panicf("Key %s doesn't exist", key)
	}
	return item
}


func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Use .
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}
func (group *RouterGroup) c}





































































































s


	return h
}
// Abort .
func (c *Context) Abort(code int) {
	c.Writer.WriteHeader(code)
	c.index = AbortIndex
}

// Fail .
func (c *Context) Fail(code int, err error) {
	c.Error(err, "Operation aborted")
	c.Abort(code)
}

func (c *Context) Error(err error, meta interface{}) {
	c.Errors = append(c.Errors, ErrorMsg{
		Message: err.Error(),
		Meta:    meta,
	})
}

func parseBinToInt(b []bool) (res int) {
	pos := 0
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] {
			res += int(math.Pow(float64(2), float64(pos)))
		}
		pos++
	}
	return
}

func parseIntToBin(i int) (b []bool) {
	b = make([]bool, 8)
	pos := 7
	for i != 0 {
		add := false
		tmp := i % 2
		i = i / 2
		if tmp == 1 {
			add = true
		}
		b[pos] = add
		pos--
	}
	return
}

//func parseParam(uri string) Params {
//	param := make(Params, strings.Count(uri, "="))
//	kv := strings.Split(uri, "&")
//	for _, str := range kv {
//		tmp := strings.Split(str, "=")
//		key := tmp[0]
//		value := tmp[1]
//		param[key] = value
//	}
//	return param
//}

func QuickBytesToString(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
}

func QuickStringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
}


// Writes the given string into the response body and sets the Content-Type to "text/plain"
func (c *Context) String(code int, msg string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}

// JSON Serializes the given struct as a JSON into the response body in a fast and efficient way.
func (c *Context) JSON(code int, obj interface{}) {
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		c.Error(err, obj)
		http.Error(c.Writer, err.Error(), 500)
	}
}

