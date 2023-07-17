package Geektutu_learning

import (
	"fmt"
	"net/http"
	"testing"
)

//func Login(c *Context) {
//	fmt.Print(c.params)
//	fmt.Print()
//	_, _ = c.w.Write([]byte("Login 请求成功\n"))
//
//}
//func FileIndex(c *Context) {
//	fmt.Print(c.params)
//	//fmt.Print()
//	_, _ = c.w.Write([]byte("文件路由" + c.params["filename"] + "\n"))
//}
//func ParamsIndex(c *Context) {
//	fmt.Print(c.params)
//	//fmt.Print()
//	_, _ = c.w.Write([]byte("参数路由" + c.params["name"] + "\n"))
//}
//func ParamsIndex2(c *Context) {
//	fmt.Print(c.params)
//	//fmt.Print()
//	_, _ = c.w.Write([]byte("参数路由2" + c.params["name"] + "\n"))
//}

// func Register(c *Context) { // 对应的处理器
//
//		_, _ = c.w.Write([]byte("Register 请求成功\n"))
//	}
func TestHTTP_Start(t *testing.T) {
	h := NewHTTP()
	//h.GET("/login/jzf", Login)
	//h.GET("/login/:name", ParamsIndex)
	//h.GET("/", Login)
	h.GET("/login/jzf", func(c *Context) { // 相应的视图函数是什么？
		c.String(http.StatusOK, fmt.Sprintf("静态路由:%s\n", c.Pattern))
	})
	//h.GET("/login/name", ParamsIndex)
	h.GET("/get/*filename", func(c *Context) {
		filename, err := c.Params("filename") // 找动态路由的方法
		if err != nil {
			c.String(http.StatusNotFound, "通配符参数错误\n")
			return
		}
		c.HTML(http.StatusOK, fmt.Sprintf("通配符路由:%s\n", filename))
	})
	h.GET("/study/:language", func(c *Context) {
		language, err := c.Params("language") // 找动态路由的方法
		if err != nil {
			c.String(http.StatusNotFound, "参数错误\n")
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("参数路由:%s\n", language))
	})

	//h.GET("/login/:filename", FileIndex)
	//h.GET("/post/:filename", ParamsIndex)
	//h.GET("/", Login)
	//h.POST("/register", Register)
	err := h.Start(":8888")
	if err != nil {
		panic(err)
	}
}
