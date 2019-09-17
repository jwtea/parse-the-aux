package main

import (
	log "github.com/sirupsen/logrus"
)

type Item struct {
	Key   string
	Value string
}

func main() {
	p := NewParser("aux-addon.lua")
	log.Print(p.GetItems())
}
