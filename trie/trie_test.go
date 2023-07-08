package trie

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.AddRouter("/user/login", "hello")
	r.AddRouter("/user/login/study", "hello")
	r.AddRouter("/user/register/jzf", "hello")
	r.AddRouter("/user/study/python", "hello")
	return r
}

//	func TestRouter_AddRouter(t *testing.T) {
//		r := newRouter()
//		testCases := []struct {
//			name       string
//			pattern    string
//			data       string
//			wantRouter *Router
//		}{
//			{
//				name:    "xxx",
//				pattern: "/user/login",
//				data:    "hello",
//				wantRouter: &Router{map[string]*Node{
//					"/": {
//						part: "/",
//						children: map[string]*Node{
//							"user": {
//								part: "user",
//								children: map[string]*Node{
//									"login": {
//										part: "login",
//										data: "hello",
//									},
//								},
//							},
//						},
//					},
//				}},
//			},
//		}
//		router := &Router{map[string]*Node{
//			"/": {
//				part: "/",
//			},
//		}}
//		for _, tc := range testCases {
//			t.Run(tc.name, func(t *testing.T) {
//				router.AddRouter(tc.pattern, tc.data)
//				assert.Equal(t, tc.wantRouter, router)
//			})
//		}
//	}
func TestRouter_GetRouter(t *testing.T) {
	r := newTestRouter()
	testCases := []struct {
		// 测试的名字，任意给就好
		name string
		// 想要匹配的节点
		findPattern string
		// 想要返回的数据
		wantData string
		// 理想中的错误
		wantErr error
	}{
		{
			name:        "success",
			findPattern: "/user/login",
			wantData:    "hello",
		},
		{
			name:        "success",
			findPattern: "/user/login/study",
			wantData:    "hello",
		},
		{
			name:        "error",
			findPattern: "/user//login",
			wantErr:     errors.New("pattern格式不对"),
		},
		{
			name:        "error2",
			findPattern: "/userasjhd/logi/n",
			wantErr:     errors.New("pattern不存在"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := r.GetRouter(tc.findPattern)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantData, n.data)
		})
	}
}
