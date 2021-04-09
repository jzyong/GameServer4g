package util

import (
	"bytes"
	"encoding/json"
	"fmt"
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
