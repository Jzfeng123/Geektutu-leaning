package Geektutu_leaning

import "testing"

func TestHTTP_Start(t *testing.T) {
	h := NewHTTP()
	if err := h.Start(":8080"); err != nil {
		t.Fail()
	}
}
