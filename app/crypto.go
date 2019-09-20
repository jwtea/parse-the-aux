package main

import (
	"hash/fnv"
)

func getHash32(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}
