package Geektutu_leaning

import "net/http"

type Context struct {
	// 相应
	w http.ResponseWriter
	// 请求
	req *http.Request
	// 请求方式
	Method string
	// URL
	Pattern string
	// 相应信息
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:       w,
		req:     r,
		Method:  r.Method,
		Pattern: r.URL.Path,
	}
}
