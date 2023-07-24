package Geektutu_learning

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 方便用户操作
type H map[string]interface{}

type Context struct {
	// 响应
	w http.ResponseWriter
	// 请求
	req *http.Request
	// 请求方式
	Method string
	// URL
	Pattern string
	//  动态路由参数
	params map[string]string
	// 请求相关的信息
	// 1.请求参数: GET /user/:id 获取ID是1的用户信息, DELETE /user/:id 删除ID是1的用户信息
	// 2.查询参数 Query http://localhost:8888?srarch=JZF&page=10 --> {search:JZF，page=10}
	// 3.请求体
	// 维护一份查询数据
	cacheQuery url.Values
	// 维护一份请求数据
	cacheBody io.ReadCloser
	// 响应相关信息
	// 状态码
	StatusCode int
	// 响应头
	header map[string]string
	// 响应体
	data []byte
}

// 获取动态路由的方法
func (c *Context) Params(key string) (string, error) {
	value, ok := c.params[key]
	if !ok { //如果找不到
		return "", errors.New(fmt.Sprintf("找不到key[%s]对应的value", key))
	}
	return value, nil
}

// 查询参数
func (c *Context) Query(key string) (string, error) {
	// 维护一份查询数据
	if c.cacheQuery == nil {
		c.cacheQuery = c.req.URL.Query() //存缓存里面
	}
	/*
		getParam, ok := c.cacheQuery.Get(key)
		如果使用这种办法来取key中的值的话，返回的值就都是字符串，这样就无法判断当返回空字符串时，是由于key找不到还是当前key不存在值的情况
		因此采用以下这种办法
	*/
	getParam, ok := c.cacheQuery[key]
	if ok { //这样就可以避免切片为空没有数据的情况
		return getParam[0], nil //这么写的话无法处理返回切片为空的情况
	}
	return "", errors.New(fmt.Sprintf("web: [%s]不存在", key))
}

// 获取请求体的数据
func (c *Context) Form(key string) (string, error) {
	if c.cacheBody == nil {
		c.cacheBody = c.req.Body
	}
	// 取数据
	err := c.req.ParseForm() //ParseForm无法取出JSON格式的数据
	// "application/x-www-form-urlencoded"只能取出这一种格式的数据, 因此得自己定义一个使之能读
	if err != nil {
		return "", err
	}
	return c.req.FormValue(key), nil
}

// 获取JSON格式的请求体数据
func (c *Context) BindJSON(dest any) error {
	if c.cacheBody == nil {
		c.cacheBody = c.req.Body
	}
	decoder := json.NewDecoder(c.cacheBody) // Decoder解码
	decoder.DisallowUnknownFields()         //处理未知字符
	return decoder.Decode(dest)
}

// 新建一个上下文，就相当于是一个视图函数
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:       w,
		req:     r,
		Method:  r.Method,
		Pattern: r.URL.Path,
		header:  map[string]string{},
	}
}

// 直接在上下文维护数据
// SetStatusCode 响应状态码,
func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	//c.w.WriteHeader(code) //直接往响应体写, 但只能写入这一次
}

func (c *Context) SetHeader(key, value string) {
	c.header[key] = value
	//c.w.Header().Set(key, value) // 直接往响应体写, 但只能写这一次
}
func (c *Context) DelHeader(key string) {
	delete(c.header, key)
	//c.w.Header().Set(key, value) // 直接往响应体写, 但只能写这一次
}
func (c *Context) SetData(data []byte) {
	c.data = data
}

// y以上三种方式是小零件，需要配合其他的方法给用户使用 --> JSON, HTML, String
// 纯文本类型
func (c *Context) String(code int, data string) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	c.SetData([]byte(data))
}

// JSON 响应JSON格式数据
func (c *Context) JSON(code int, data any) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(code)
	res, err := json.Marshal(data) //序列化处理
	if err != nil {
		// 1.存在的问题：如果panic报错
		// 我们之前设置的响应头和状态码需要去掉吗？最好是去掉
		c.SetStatusCode(http.StatusNotFound)
		c.DelHeader("Content-Type")
		panic(err)
	}
	c.SetData(res)
}

// HTML 响应HTML
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	//c.SetHeader("", "text/html")
	c.SetStatusCode(code)
	c.SetData([]byte(html))
}

// FlashToHeader 既然我们没有直接把响应数据写入Header中，那么我们就需要自己建立一个方法将数据读进去
// 写入的顺序为响应头-状态码-响应体，顺序错的话则不能正确输出
// 这是golang的一个坑
func (c *Context) FlashToHeader() {
	// 写入响应头
	for key, val := range c.header {
		c.w.Header().Set(key, val) //优先级最低所以放最上面先去执行
	}
	// 写入状态码
	c.w.WriteHeader(c.StatusCode)
	// 写入响应体, 写给客户端
	_, _ = c.w.Write(c.data)
}

//我们在Context上下文中，维护一些请求相关的数据是因为可能有很多视图函数需要用到这些数据
//那为什么还需要维护一些响应数据呢？
// 由于直接使用wirte.Setheader这种方式具有一定的局限性，不利于修改和维护。
