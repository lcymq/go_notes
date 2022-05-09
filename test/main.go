package main

import (
	"fmt"
	"strings"
)

func main() {
	s1 := "abcdef"
	s2 := "bc"
	res := strings.Split(s1, s2)
	fmt.Println(res)
}

func Split(str string, sep string) []string {
	res := []string{}
	temp := ""
	for i := range str {
		if string(str[i]) == sep {
			res = append(res, temp)
			temp = ""
		} else {
			temp += string(str[i])
		}
	}
	if temp != "" {
		res = append(res, temp)
	}
	return res
}
