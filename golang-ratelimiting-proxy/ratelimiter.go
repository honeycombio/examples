package main

import (
	"math/rand"

	"github.com/honeycombio/leakybucket"
)

type RandomRateLimiter struct {
	Frequency int // how often to reply with rate limited
}

func (r *RandomRateLimiter) Add(key string) error {
	if rand.Intn(r.Frequency) == 0 {
		return &leakybucket.BucketOverflow{}
	}
	return nil
}
func (r *RandomRateLimiter) Start() error {
	return nil
}
func (r *RandomRateLimiter) Stop() error {
	return nil
}
