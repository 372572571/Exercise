package doc

import (
	"fmt"
	"reflect"
)

func Values(obj interface{}) {
	s := reflect.ValueOf(obj).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		// 0: Name string = Sparrow
		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(),
			f.Interface())
	}
}
