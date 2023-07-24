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
func (r *router) addRouter(method string, pattern string, handler HandleFunc, middlewareChains ...MiddlewareHandleFunc) {
	fmt.Printf("add router %s - %s\n", method, pattern) // /login/jzf
	//fmt.Printf("pattern is %s \n", pattern)
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
		return // 直接退出这个func
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
		root, ok = root.addNode(part) // 循环结束之后，root是最后一个叶子节点
		if !ok {                      //返回false表示路由注册时有冲突
			panic(fmt.Sprintf("web: 路由注册冲突, %s", part))
		}
	}
	// 解决/login/name  /login/name定义了两种相同的路由的路由冲突情况
	if root.handleFunc != nil { //如果当前子节点的视图函数不为空，代表发生了路由冲突
		panic(fmt.Sprintf("web: 路由冲突(相同的注册路由) - %s", pattern))
	}
	// 设置路由相应的视图函数
	root.handleFunc = handler //给最后一个叶子节点添加上相应的视图函数
	// 设置中间件
	root.middlewareChain = middlewareChains
}

// method 不需要考虑， method直接找不到就行
// pattern可以校验一些简单的
func (r *router) getRouter(method string, pattern string) (*node, map[string]string, bool) { //当我们浏览器中直接传入的地址的时候，
	// 走的是这一个流程而不需要重新注册一遍路由
	// 我们注册的路由是/study/:course， 但用户传进来的路由是/study/golang,应该考虑如何匹配:course和golang的问题，用一个map来表示匹配关系。
	params := map[string]string{} // 两个都是字符串
	if pattern == "" {
		return nil, params, false
	}
	root, ok := r.trees[method]
	if !ok { // 如果没有这个方法的节点就直接返回false不需要再添加，这是查询，不是注册，复制粘贴都能错- -||
		//r.trees[method] = &node{
		//	part: "/",
		//}
		//root = r.trees[method]
		return nil, params, false
	}
	// 如果传入的是"/",直接返回对应的root就行
	if pattern == "/" {
		return root, params, true
	}
	// /user/login/ --> 这种是合理的，因此应该考虑将开头结尾的/去掉,
	//包括/login/jzf///////这种情况也是合理的，因为strings.Trim(pattern, "/")会将前后所有/都给去掉
	parts := strings.Split(strings.Trim(pattern, "/"), "/") //
	for _, part := range parts {
		if part == "" {
			return nil, params, false
		}
		root = root.getNode(part) // 获取路由
		// 由于我想打印找到后的root是什么样子，但是root是一个空指针形式，不能直接使用format打印出来，因此会报错
		//_, err := fmt.Printf("root is %s\n", root.part) // 这种方法是会返回一个err的
		//if err != nil {
		//	panic(err)
		//}
		//print("root is", root) //这种打印能打印空指针就不会报错，但是无法打印出指针指向的元素
		if root == nil { //没有注册这个路由的话返回空值直接return
			return nil, params, false
		}
		// {course: golang}，参数路由的特殊处理
		if strings.HasPrefix(root.part, ":") { //参数路由
			params[root.part[1:]] = part
		}
		// /login/*filename/jzf/temple.css --> {filename:/jzf/temple.css}
		// 从pattern里面找
		// 也是存在params里面
		if strings.HasPrefix(root.part, "*") { //通配符匹配，贪婪匹配
			// 找到filename所对应的路由匹配下标
			index := strings.Index(pattern, part) //找到当前part在pattern中的位置，首字母开头的位置
			// 从pattern里面找
			fmt.Print(index)
			params[root.part[1:]] = pattern[index:]
			return root, params, root.handleFunc != nil // 通配符节点一定是叶子节点。
		}
		// 当出现/study/:course/golang时候，应该让他继续匹配
	}
	return root, params, root.handleFunc != nil //root.handleFunc != nil判断是否是最后一个节点
}

// 构造前缀树节点
type node struct {
	part string
	// 子节点，
	children map[string]*node
	// 处理器-视图函数
	handleFunc HandleFunc
	// 中间件列表
	middlewareChain MiddlewareChains //这玩意是由用户传进来的，在某个视图函数上所规定的中间件
	// 参数路由
	// 为什么是一个纯的node节点，因为动态路由是变化的，不好去判断当前节点属于哪一个参数
	// 静态路由和动态路由的优先级问题 ---> 静态路由优先级高于动态路由
	paramChild *node
	//通配符匹配 贪婪匹配
	starChild *node
}

// addNode 在服务启动前调用
func (n *node) addNode(part string) (*node, bool) {
	/* 有一个问题就是
	/login/:name/jzf 最后一个jzf的节点会带一个视图函数
	/login/:name    :name会带一个视图函数，两者在前缀树上属于同一前缀,怎么去处理
	*/
	// 因此不能直接粗暴的如以下设计，这会导致注册第二个:name路由的时候无法return回当前这个参数点
	//if strings.HasPrefix(part, ":") && n.paramChild == nil {n.paraChild = &node{part:part};return n.paraChild}
	// 实现参路路由和贪婪匹配
	if strings.HasPrefix(part, "*") { //通配符路由
		if n.paramChild != nil { // 不允许 /login/:name 和 /login/*filename这种情况出现，这是路由冲突
			return nil, false
		}
		if n.starChild == nil {
			n.starChild = &node{part: part}
		}
		return n.starChild, true
	}
	if strings.HasPrefix(part, ":") { //参数路由
		// 同样的，不允许/login/*name 和 /login/:filename这种情况出现
		if n.starChild != nil {
			return nil, false
		}
		if n.paramChild == nil { //创建参数路由
			n.paramChild = &node{part: part}
		}
		if n.paramChild.part != part { //判定 /login/:name 和/login/:filename这样的路由冲突
			return nil, false
		}
		return n.paramChild, true
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
	return child, true
}

func (n *node) getNode(part string) *node {
	// n不存在children属性
	if n.children == nil {
		if n.paramChild != nil { // 一个路由的同一个位置不能同时有静态路由和动态路由
			return n.paramChild
		}
		if n.starChild != nil {
			return n.starChild
		}
		return nil
	}
	child, ok := n.children[part]
	// 静态路由优先级高于动态路由
	if !ok {
		// 到这里表示没有匹配到静态路由，因此可以考虑是否存在动态路由
		if n.paramChild != nil {
			return n.paramChild
		}
		if n.starChild != nil {
			return n.starChild
		}
		return nil
	}
	return child
}

/**
	路由冲突存在的情况,三种
	/study/login
	/study/login
	-------------
	/study/:course
	/study/:action
	-------------
	/study/*course
	/study/:course
**/

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

/* 添加节点的操作
如果能够添加成参数路由，那一定能够添加成静态路由，所以应该先从动态路由开始判断一直到静态路由。

*/
