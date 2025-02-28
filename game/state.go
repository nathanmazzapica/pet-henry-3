package game

import "sync/atomic"

var Counter int64

func IncrementPetCounter() {
	atomic.AddInt64(&Counter, 1)
}
