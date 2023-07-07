package Geektutu_leaning

import (
	"testing"
)

func Login(c *Context) {
	c.w.Write([]byte("Login 请求成功\n"))
}

func Register(c *Context) {
	c.w.Write([]byte("Register 请求成功\n"))
}
func TestHTTP_Start(t *testing.T) {
	h := NewHTTP()
	h.GET("/login", Login)
	h.POST("/register", Register)
	err := h.Start(":8888")
	if err != nil {
		panic(err)
	}
}
