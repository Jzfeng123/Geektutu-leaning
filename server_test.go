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
	h.Use(Logger())
	//v1.Use(Logger())
	//h.GET("/login/jzf", Login)
	//h.GET("/login/:name", ParamsIndex)
	//h.GET("/", Login)
	h.GET("/login/jzf", func(c *Context) { // 相应的视图函数是什么？
		//fmt.Println("123", h.prefix)
		c.String(http.StatusOK, fmt.Sprintf("静态路由:%s\n", c.Pattern))
		//}, func(next HandleFunc) HandleFunc {
		//	return func(c *Context) {
		//		fmt.Println("hello i am coming1")
		//		next(c)
		//		fmt.Println("hello i am go1")
		//	}
		//}, func(next HandleFunc) HandleFunc {
		//	return func(c *Context) {
		//		fmt.Println("hello i am coming2")
		//		next(c)
		//		fmt.Println("hello i am go2")
		//	}
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
	h.GET("/login/:language", func(c *Context) {
		language, err := c.Params("language") // 找动态路由的方法
		if err != nil {
			c.String(http.StatusNotFound, "参数错误\n")
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("参数路由:%s\n", language))
	})
	//h.GET("/json", func(c *Context) {
	//	c.JSON(http.StatusOK, H{ //map[string]interface //这里调用的JSON保存的格式是map[string]string
	//		"code": 200,
	//		"msg":  "请求成功",
	//		"data": []string{
	//			"A", "B", "C",
	//		},
	//	})
	//})
	//// 测试HTML
	//h.GET("/html", func(c *Context) {
	//	c.HTML(http.StatusOK, `<h1 style="color:red;">hello world</h1>`)
	//	//``不会解析字符串的转义符，""会解析字符串的转义符
	//	//效率是前者高于后者
	//})
	//
	//h.GET("/query", func(c *Context) {
	//	username, err := c.Query("username")
	//	if err != nil {
	//		c.SetStatusCode(http.StatusNotFound)
	//		return
	//	}
	//	passwd, err := c.Query("passwd")
	//	if err != nil {
	//		c.SetStatusCode(http.StatusNotFound)
	//		return
	//	}
	//	c.JSON(http.StatusOK, H{
	//		"code":     200,
	//		"msg":      "请求成功",
	//		"username": username,
	//		"passwd":   passwd,
	//	})
	//})
	////h.Use(Logger())
	//h.POST("/body", func(c *Context) {
	//	type User struct {
	//		Username string `json:"username"`
	//		Passwd   string `json:"passwd"`
	//	}
	//	var user User
	//	err := c.BindJSON(&user)
	//	if err != nil {
	//		c.String(http.StatusNotFound, err.Error())
	//		return
	//	}
	//	//username, err := c.Form("username")
	//	//if err != nil {
	//	//	c.SetStatusCode(http.StatusNotFound)
	//	//	return
	//	//}
	//	//passwd, err := c.Form("passwd")
	//	//if err != nil {
	//	//	c.SetStatusCode(http.StatusNotFound)
	//	//	return
	//	//}
	//	c.JSON(http.StatusOK, H{
	//		"code":     http.StatusOK,
	//		"msg":      "请求成功",
	//		"username": user.Username,
	//		"passwd":   user.Passwd,
	//	})
	//})
	//h.GET("/login/:filename", FileIndex)
	//h.GET("/post/:filename", ParamsIndex)
	//h.GET("/", Login)
	//h.POST("/register", Register)
	//// 假设用户栾川
	//v1 := h.Group("//v1//")
	//v1.Use(Logger())
	//{
	//	v1.GET("/login/xxx", func(c *Context) {
	//
	//		c.HTML(http.StatusOK, `<h1 style="color:red;">XX monster</h1>`)
	//	})
	//	v1.GET("/login/:name", func(c *Context) {
	//		c.HTML(http.StatusOK, `<h1 style="color:blue;">XX guai</h1>`)
	//	}, func(next HandleFunc) HandleFunc {
	//		return func(c *Context) { //某个视图函数上的中间件测试
	//			fmt.Println("我是1号我来了")
	//			next(c)
	//			fmt.Println("我是1号我润了 ")
	//		}
	//	}, func(next HandleFunc) HandleFunc {
	//		return func(c *Context) {
	//			fmt.Println("我是2号我来了")
	//			next(c)
	//			fmt.Println("我是2号我润了")
	//		}
	//	})
	//	v1.GET("/acquire/*filename", func(c *Context) {
	//		filename, err := c.Params("filename") // 找动态路由的方法
	//		if err != nil {
	//			c.String(http.StatusNotFound, "没这玩意")
	//			return
	//		}
	//		c.JSON(http.StatusOK, H{
	//			"code":     200,
	//			"msg":      "请求成功",
	//			"name":     "hello jzf",
	//			"filename": filename,
	//		})
	//	})
	//}
	////v1下面所有的路由都会走Logger
	//v2 := v1.Group("v2//")
	//{
	//	v2.GET("/login/xx", func(c *Context) {
	//		c.HTML(http.StatusOK, `<h1 style="color:red;">Hello World</h1>`)
	//	})
	//	v2.GET("/login/:name", func(c *Context) {
	//		c.HTML(http.StatusOK, `<h1 style="color:blue;">Whoo!!</h1>`)
	//	})
	//	v2.GET("/acquire/*filename", func(c *Context) {
	//		filename, err := c.Params("filename") // 找动态路由的方法
	//		if err != nil {
	//			c.String(http.StatusNotFound, "没这玩意")
	//			return
	//		}
	//		c.JSON(http.StatusOK, H{
	//			"code":     200,
	//			"msg":      "请求成功",
	//			"name":     "king",
	//			"filename": filename,
	//		})
	//	})
	//}
	err := h.Start(":8888")
	if err != nil {
		panic(err)
	}
}
