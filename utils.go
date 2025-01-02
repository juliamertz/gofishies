package main

import (
	"math/rand/v2"
	"time"
)

func RandOneIn(n int) bool {
	src := rand.NewPCG(uint64(n), uint64(time.Now().UnixNano()))
	rng := rand.New(src)
  j := rng.IntN(n)
  return j == 1
}
