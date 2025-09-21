package array

import (
	"math/rand"
	"sync"
	"time"
)

var (
	r     = rand.New(rand.NewSource(time.Now().UnixNano()))
	rLock sync.Mutex
)

// RandomArray 随机获取数组中的一个元素
// 参数：
//   - arr: 数组
//
// 返回值：
//   - T: 数组中的一个元素
func RandomArray[T comparable](arr []T) T {
	rLock.Lock()
	defer rLock.Unlock()
	return arr[r.Intn(len(arr))]
}

// IsArrayContains 判断数组中是否包含指定元素
// 参数：
//   - arr: 数组
//   - value: 指定元素
//
// 返回值：
//   - bool: 是否包含指定元素，true表示包含，false表示不包含
func IsArrayContains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
