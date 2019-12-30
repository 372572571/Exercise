package regexpfilter

import (
	"regexp"
	"strconv"
)
// easy match function


// match [0-9a-zA-Z]{len}  [:alnum:]{len}
// fail  return false
func MatchAlnum(str string)bool{
	expr,err:=regexp.Compile("[0-9A-Za-z]{"+strconv.Itoa(len(str)-1)+"}")
	if err!=nil{
		return false
	}
	if !expr.MatchString(str){
		return false
	}
	return true
}
