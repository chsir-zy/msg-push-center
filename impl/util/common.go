package util

import "os"

/*
 * 判断文件是否存在 不存在就创建
 */
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err == nil {
			return true
		}
		return false
	}
	return true
}
