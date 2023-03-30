// @author Zongsheng Xu 2023/3/30 21:13

package utils

import (
	"math/rand"
	"sync"
	"time"
)

const (
	charset  = "0123456789"
	idLength = 16
)

var (
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu         sync.Mutex
	usedIDs    = make(map[string]bool)
)

func generateID() string {
	id := make([]byte, idLength)
	for i := range id {
		id[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(id)
}

func GenerateUniqueID() string {
	mu.Lock()
	defer mu.Unlock()

	for {
		id := generateID()
		if !usedIDs[id] {
			usedIDs[id] = true
			return id
		}
	}
}
