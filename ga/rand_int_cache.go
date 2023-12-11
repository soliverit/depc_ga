package ga

import (
	"math/rand"
)

/*
Random Cache: Quick random numbers for GAs!!!
*/
type RandIntCache struct {
	length  int
	values  []int
	counter int
}

func CreateRandIntCache(seed, max, limit int) RandIntCache {
	var randIntCache RandIntCache
	randIntCache.counter = -1 // -1 so that Next() starts on 0 [ it does increment first so it doesn't define an extra variable
	randIntCache.length = limit
	randIntCache.values = make([]int, randIntCache.length)
	r := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < limit; i++ {
		randIntCache.values[i] = r.Intn(max - 1)
	}

	return randIntCache
}
func (randIntCache *RandIntCache) Next() int {
	randIntCache.counter++
	if randIntCache.counter == randIntCache.length {
		randIntCache.counter = 0
	}
	return randIntCache.values[randIntCache.counter]
}
