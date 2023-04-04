/*
 * package utils 工具包，提供了一些常用的工具函数
 */
package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	// 使用时间戳作为随机数生成器的源种子
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 互斥锁，用于多个goroutine之间同步访问usedIDs变量
	mu sync.Mutex
	// 存储已经使用过的ID
	usedIDs = make(map[int]bool)
)

// 生成一个随机字符串作为ID
func generateID() int {

	min := int(1e15)
	max := int(1e16 - 1)
	return min + seededRand.Intn(max-min+1)

}

// 生成唯一的ID
func GenerateUniqueID() int {
	mu.Lock()
	defer mu.Unlock()

	for {
		id := generateID()
		// 如果ID没有被使用，则将其标记为已使用并返回
		if !usedIDs[id] {
			usedIDs[id] = true
			return id
		}
	}
}
