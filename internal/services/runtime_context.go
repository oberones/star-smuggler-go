package services

import (
	"math/rand"
	"time"
)

type RNG interface {
	Intn(n int) int
	Float64() float64
}

type RuntimeContext struct {
	RNG RNG
	Now func() time.Time
}

func NewRuntimeContext(rng RNG, now func() time.Time) RuntimeContext {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if now == nil {
		now = time.Now
	}

	return RuntimeContext{
		RNG: rng,
		Now: now,
	}
}

func NewSeededRuntimeContext(seed int64) RuntimeContext {
	return NewRuntimeContext(rand.New(rand.NewSource(seed)), time.Now)
}
