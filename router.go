package Geektutu_learning

import (
	"fmt"
	"strings"
)

// 根据trie构造路由森林
/*
router 结构体
trees means as follows:
{
	"GET": node{}
	"PST": node{}
	"DELETE": node{}
	...
}
*/
type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

// addRouter 注册路由
// 需要考虑什么样的pattern是合法的？
// 不需要考虑method的原因是，method传入是固定的请求方法，就算是乱传那就到时候返回nil就好
// 我们可以规定用户必须以/开头并且不以/结尾,将错误panic给用户
// 还需要考虑的是pattern中是否会有长连/的情况，例如/user//jzf///login/////
// 还需要考虑 / 是一个正规的路由
func (r *router) addRouter(method string, pattern string, handler HandleFunc) {
	fmt.Printf("add router %s - %s\n", method, pattern) // /login/jzf
	fmt.Printf("pattern is %s \n", pattern)
	if pattern == "" { //
		panic("路由不可以为空")
	}
	// TODO 如果根路由是/怎么办？
	root, ok := r.trees[method] //root -> *node
	if !ok {                    //根节点不存在，创一个
		root = &node{
			part: "/", // 默认的根节点
		}
		r.trees[method] = root
	}
	if pattern == "/" {
		root.handleFunc = handler
		return // 直接宕掉这个func
	}
	if !strings.HasPrefix(pattern, "/") {
		panic("路由必须以 / 开头")
	}
	if strings.HasSuffix(pattern, "/") {
		panic("路由不能以 / 结尾")
	}
	//switch {
	//case pattern == "":
	//	panic("路由不可以为空\n")
	//case !strings.HasPrefix(pattern, "/"):
	//	panic("路由必须以 / 开头\n")
	//case strings.HasSuffix(pattern, "/"):
	//	panic("路由不能以 / 结尾\n")
	//}
	// 获取根节点
	// 第一版写法：
	//root, ok := r.trees[method] //root -> *node
	//if !ok {                    //根节点不存在，创一个
	//	root = &node{
	//		part: "/", // 默认的根节点
	//	}
	//	r.trees[method] = root
	//}
	// 切割pattern
	parts := strings.Split(pattern[1:], "/")
	for _, part := range parts {
		if part == "" {
			panic("web路由不能连续出现 / \n")
		}
		root = root.addNode(part) // 循环结束之后，root是最后一个叶子节点
	}
	root.handleFunc = handler //给最后一个叶子节点添加上相应的视图函数
}

// method 不需要考虑， method直接找不到就行
// pattern可以校验一些简单的
func (r *router) getRouter(method string, pattern string) (*node, bool) { //当我们浏览器中直接传入的地址的时候，
	// 走的是这一个流程而不需要重新注册一遍路由
	if pattern == "" {
		return nil, false
	}
	root, ok := r.trees[method]
	if !ok {
		r.trees[method] = &node{
			part: "/",
		}
		root = r.trees[method]
	}
	// 如果传入的是"/",直接返回对应的root就行
	if pattern == "/" {
		return root, true
	}
	// /user/login/ --> 这种是合理的，因此应该考虑将开头结尾的/去掉,
	//包括/login/jzf///////这种情况也是合理的，因为strings.Trim(pattern, "/")会将前后所有/都给去掉
	parts := strings.Split(strings.Trim(pattern, "/"), "/") //
	for _, part := range parts {
		if part == "" {
			return nil, false
		}
		root = root.getNode(part)
		if root == nil {
			return nil, false
		}
	}
	return root, true
}

// 构造前缀树节点
type node struct {
	part string
	// 子节点，
	children map[string]*node
	// 处理器-视图函数
	handleFunc HandleFunc
	// 参数路由
	// 为什么是一个纯的node节点，因为动态路由是变化的，不好去判断当前节点属于哪一个参数
	// 静态路由和动态路由的优先级问题
	paramChild *node
}

// 在服务启动前调用
func (n *node) addNode(part string) *node {
	if strings.HasPrefix(part, ":") && n.paramChild == nil { //参数路由
		n.paramChild = &node{part: part}
		return n.paramChild
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[part]
	if !ok { //如果当前节点没有part这一个属性，那么就造一个
		child = &node{
			part: part,
		}
		//
		n.children[part] = child
	}
	return child
}

func (n *node) getNode(part string) *node {
	if n.children == nil {
		return nil
	}
	child, ok := n.children[part]
	if !ok {
		return nil
	}
	return child
}

/*
路由分为动态和静态路由：
- 静态路由 /user/login/   /study/golang --> 规定好的
- 动态路由:
  1. 参数路由
     /study/:course 咱们注册的路由，匹配的时候可能会匹配到/study/golang、 study/python，
					但是/study/golang/action这种路由匹配不到
  2. 通配符路由: 贪婪匹配
	/static/*filepath 这是咱们注册的路由
		匹配的时候可能匹配到/static/css/style.css
						/static/js/index.js
  3. 正则路由
*/
