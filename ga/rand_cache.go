package ga

import (
	"math/rand"
)

/*
Random Cache: Quick random numbers for GAs!!!
*/
type RandCache struct {
	length  int
	values  []float32
	counter int
}

func CreateRandCache(seed, limit int) RandCache {
	var randCache RandCache
	randCache.counter = -1 // -1 so that Next() starts on 0 [ it does increment first so it doesn't define an extra variable
	randCache.length = limit
	randCache.values = make([]float32, randCache.length)
	r := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < limit; i++ {
		randCache.values[i] = r.Float32()
	}

	return randCache
}
func (randCache *RandCache) Next() float32 {
	randCache.counter++
	if randCache.counter == randCache.length {
		randCache.counter = 0
	}
	return randCache.values[randCache.counter]
}
