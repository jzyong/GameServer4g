package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

//ToStringIndent 将任意结构转化为json缩进后的字符串 方便输出调试
func ToStringIndent(what interface{}) string {
	b, err := json.Marshal(what)
	if err != nil {
		return fmt.Sprintf("%+v", what)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	return out.String()
}

//ToString 将任意结构转化为json字符串 方便输出调试
func ToString(what interface{}) string {
	b, err := json.Marshal(what)
	if err != nil {
		return fmt.Sprintf("%+v", what)
	}
	return string(b)
}

//是否包含某个元素
func SliceContains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}
