package main

import (
	"fmt"
)

func main() {
	p := NewParser("aux-addon.lua")

	items := p.GetItems()
	for i, wi := range items {
		fmt.Printf("i:%d time:%v cValue:%d \n", i, wi.Timestamp, wi.MinBuyout)
	}
}
