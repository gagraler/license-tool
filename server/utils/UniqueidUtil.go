/*
 * package utils 工具包，提供了一些常用的工具函数
 */
package utils

import (
	"math/rand"
	"sync"
	"time"
)

const (
	// 字符集为数字0-9
	charset = "0123456789"
	// ID长度为16位
	idLength = 16
)

var (
	// 使用时间戳作为随机数生成器的源种子
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 互斥锁，用于多个goroutine之间同步访问usedIDs变量
	mu sync.Mutex
	// 存储已经使用过的ID
	usedIDs = make(map[string]bool)
)

// 生成一个随机字符串作为ID
func generateID() string {
	id := make([]byte, idLength)
	for i := range id {
		id[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(id)
}

// 生成唯一的ID
func GenerateUniqueID() string {
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
