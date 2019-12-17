package utils

// ArrayIndexOf ...
// 判断当前数组是否存在这个下标
// 存在返回true
func ArrayIndexOf(array []interface{}, index int) bool {
	var l = len(array)
	for i := 0; i < l; i++ {
		if i == index {
			return true
		}
	}
	return false
}
