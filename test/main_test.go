package main

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	got := Split("a:b:", ":")
	want := []string{"a", "b"}
	if !reflect.DeepEqual(want, got) { // 切片是引用类型，不能直接比较
		t.Errorf("expected:%v, got:%v", want, got)
	} else {
		t.Log(got)
	}
	t.Run()
}
