package common

/*********************************** 公共函数 *************************************/

//检查一个string类型的值在slice中是否存在
func IsStringExistsInSlice(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}


