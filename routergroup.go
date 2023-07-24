package Geektutu_learning

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
RouterGroup 里面需要什么属性
// 需要一个唯一标识
*/
//
type RouterGroup struct {
	// 唯一标识, 先考虑这一个信息
	prefix string
	// 继承engine，由于我们将method方法移到了路由组里面，我们需要使用engine中的addRoute方法。
	engine *HTTPServer
	// 父节点表示父路由组 实现路由组的嵌套例如/v1/v2/login/jzf /v1/v2/study/golang
	// v2和v1都是路由组，但是v1是v2的父路由组,因此我们可以先注册v1的路由组，然后v2在v1的基础上注册路由，此时parent为v1就可以达到我们要的效果。
	parent *RouterGroup
	// middlewares 当前路由组所有中间件
	middlewares []MiddlewareHandleFunc // 责任链机制中的多层嵌套
}

// Group 注册路由组 之后应该如何将注册后的路由组添加到框架之中？
func (group *RouterGroup) Group(prefix string) *RouterGroup { // 嵌套的时候group里面是有前一个路由组的信息的
	// 我们应该对传入的prefix进行校验工作防止用户乱传
	// 类似这样的瞎搞 "/v1/" "v1" "v1/"
	//prefix = strings.Trim(prefix, "/")
	//prefix = "/" + group.prefix + prefix //默认给他变成"/v1"的形式
	prefix = fmt.Sprintf("/%s", strings.Trim(prefix, "/"))
	rg := &RouterGroup{ //没有绑定engine导致注册失败
		prefix: group.prefix + prefix,
		engine: group.engine,
		parent: group, // group赋值给父节点
	}
	// 将新建的路由组添加到HTTPServer
	group.engine.groups = append(group.engine.groups, rg)
	return rg
}

// Use 首先我们需要注册中间件-->之后我们得匹配中间件，那么就得到ServerHTTP中
// 将中间件应用到某个group上
func (group *RouterGroup) Use(mids ...MiddlewareHandleFunc) {
	// 问题: 中间件应该放在哪里？ -->  应该在每个路由组上都维护一个中间件
	group.middlewares = append(group.middlewares, mids...) //维护一个中间件组
}

// NewGroup 抽取出的公共方法
// 注册完路由组之后，我们应该把注册后的路由组添加到框架当中，我们的指令是
/*
	v1 := h.group("/v1") //目前只实现了这一个,prefix
	v1.GET("/login/study") //应该如何实现这个，那么就是在对应的方法下将prefix添加到pattern的头部，URL就变成了/v1/login/study
*/
func NewGroup() *RouterGroup { //给每个结构体加一个构造方法为了方便扩展
	return &RouterGroup{}
}

// 这样就可以把路由组添加进框架之中,同时这里也是注册路由的唯一路径
func (group *RouterGroup) addRouter(method string, part string, handleFunc HandleFunc, middlewareChains ...MiddlewareHandleFunc) {
	pattern := group.prefix + part //合并prefix和part变成一个新的pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handleFunc, middlewareChains...)
}

// GET
func (group *RouterGroup) GET(pattern string, handleFunc HandleFunc, middlewareChains ...MiddlewareHandleFunc) {
	//我们可以这样子来进行一个添加路由的操作
	//pattern = group.prefix + pattern, 但是这样的操作太过于冗余，这样就得每一种方法都得添加一次这样的替换，因此我们可以将它封装成一个方法。
	group.addRouter(http.MethodGet, pattern, handleFunc, middlewareChains...)
}

// POST
func (group *RouterGroup) POST(pattern string, handleFunc HandleFunc, middlewareChains ...MiddlewareHandleFunc) {
	group.addRouter(http.MethodPost, pattern, handleFunc, middlewareChains...)
}

// PUT
func (group *RouterGroup) PUT(pattern string, handleFunc HandleFunc, middlewareChains ...MiddlewareHandleFunc) {
	group.addRouter(http.MethodPut, pattern, handleFunc, middlewareChains...)
}

// DELETE
func (group *RouterGroup) DELETE(pattern string, handleFunc HandleFunc, middlewareChains ...MiddlewareHandleFunc) {
	group.addRouter(http.MethodDelete, pattern, handleFunc, middlewareChains...)
}

// 每一个路由组都需要维护一个HTTPServer
// 我们注册一个或多个中间件到某一个具体的视图函数上，那这些注册的中间件需要维护在哪里？
// 当然应该是维护在当前路由的节点上 --> 因此我们需要在router中加入中间件
/*
相比较于gin框架的洋葱模型，责任链模型不方便于让用户添加全局中间件，这是由于我们这一块将中间件函数与视图函数看成了是两个不同的东西，因此每一次
注册的时候传入的方法需要绑定上相应的中间件函数有就是有没有就是没有，并不能做到一起赋值的目的，除非你是个路由组。采用责任链模式如果要加上全局中间件
的话，则需要在ServerHTTP里面定义相应的全局中间件，否则就只能一个个声明
但是Gin框架的洋葱模型是一种自动添加中间件的机制，只需要调用r.Use方法就可以将相应的中间件调用为全局中间件或局部中间件。
*/
