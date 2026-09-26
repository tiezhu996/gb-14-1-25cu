// Package strutil 提供无业务依赖的通用字符串工具（可复用库放 pkg/）。
package strutil

// Contains 判断字符串切片是否包含目标值。
func Contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
