package genetic_algorithm

import (
	"math/rand"
	"time"
)

func shuffleTwo(a, b int) (int, int) {
	c := getDiscreteRandom(0, 1)
	if c == 0 {
		return a, b
	} else {
		return b, a
	}
}

func getContinuousRandom(min, max float64) float64 {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

func getDiscreteRandom(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func getDiscreteRandomExcludeA(min, max, a int) int {
	var b int
	for i := 0; i < 1; {
		b = getDiscreteRandom(min, max)
		if b != a {
			break
		}
	}
	return b
}

func getCategoricalRandom(category []string) int {
	return rand.Intn(len(category))
}

func getCategoricalRandomExcludeA(category []string, a string) int {
	var b int
	for i := 0; i < 1; {
		b = getCategoricalRandom(category)
		if category[b] != a {
			break
		}
	}
	return b
}
