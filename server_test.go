package Geektutu_learning

import (
	"testing"
)

func Login(c *Context) {
	_, _ = c.w.Write([]byte("Login 请求成功\n"))
}

func Register(c *Context) { // 对应的处理器
	_, _ = c.w.Write([]byte("Register 请求成功\n"))
}
func TestHTTP_Start(t *testing.T) {
	h := NewHTTP()
	h.GET("/login/jzf", Login)
	h.GET("/login", Login)
	h.GET("/", Login)
	//h.GET("/", Login)
	h.POST("/register", Register)
	err := h.Start(":8888")
	if err != nil {
		panic(err)
	}
}
