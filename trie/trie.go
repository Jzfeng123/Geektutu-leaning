package trie

import (
	"errors"
	"strings"
)

// 实现一个简单的前缀树路由，假设用户传入的URL的标准的
type router struct {
	roots map[string]*node
}

// 初始化一个router
func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
	}
}

// /user/login
// /user/register
// 1.首先默认用户传的数据是标准的，如上
func (r *router) AddRouter(pattern string, data string) {
	_, ok := r.roots["/"] // 根路由"/"
	if !ok {              //创建根路由
		r.roots["/"] = &node{}
	}
	// 这一段可以抽象到另一个函数上去
	// 将pattern中的URL根据/进行划分形成切片
	root := r.roots["/"]                                    // "/"后期为method
	parts := strings.Split(strings.Trim(pattern, "/"), "/") //去除左右两边及中间的/
	for _, part := range parts {
		if part == "" {
			panic("pattern不符合格式")
		}
		root = root.Insert(part)
	}
	// 循环结束之后，我们的root会来到叶子节点
	root.data = data // 此时设置data的值
}

// 获取路由
func (r *router) GetRouter(pattern string) (*node, error) {
	root, ok := r.roots["/"]
	if !ok {
		panic("该节点不存在")
	}
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" { // --> /user//login/jzf --> [user,  , login, jzf]，多一个空字符不允许
			return nil, errors.New("pattern格式不对")
		}
		root = root.Search(part) //调用前缀树中的查找
		if root == nil {         //如果找不到
			return nil, errors.New("pattern不存在")
		}
	}
	return root, nil
}

type node struct {
	part string // 当前节点的唯一标识
	// 用map结构保存子节点
	childern map[string]*node // 保存子节点 []*node切片slice保存
	//isWild   bool             //判断是否是模糊匹配，true代表模糊也就是通配符:
	data string //需要保存的数据
}

// 这个节点需要的功能: 注册添加和查找
func (n *node) Insert(part string) *node { //将传进来的每一个part添加到前缀树里面
	if n.childern == nil { // 如果当前n的子节点为空
		n.childern = make(map[string]*node)
	}
	child, ok := n.childern[part]
	if !ok { //如果当前节点没有part这一个属性，那么就造一个
		child = &node{
			part: part,
		}
		//
		n.childern[part] = child
	}
	return child
}
func (n *node) Search(part string) *node {
	// 如果一开始n的属性都不存在，那么就直接return nil
	if n.childern == nil {
		return nil
	}
	child, ok := n.childern[part]
	if !ok {
		return nil
	}
	return child
}
