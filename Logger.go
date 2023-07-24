package Geektutu_learning

import (
	"log"
	"time"
)

func Logger() MiddlewareHandleFunc {
	return func(next HandleFunc) HandleFunc {
		return func(c *Context) { // type HandleFunc func(*Context)
			// Start timer
			t := time.Now()
			// Process request
			log.Printf("请求进来的时间: %v", t.Format("2006-01-02 15:04:05"))
			time.Sleep(time.Second * 3)
			next(c)
			// Calculate resolution time
			log.Printf("请求的总时间: %v", time.Now().Format("2006-01-02 15:04:05"))
			log.Printf("[%d] %s in %v", c.StatusCode, c.Pattern, time.Since(t))
		}
	}
}
