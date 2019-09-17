package main

import log "github.com/sirupsen/logrus"

func main() {
	p := NewParser("aux-addon.lua")

	items := p.GetItems()

	so := NewStoreOptions()
	store := NewStore(so)

	for _, wi := range items {
		log.Infof("Item: %v", wi)
		store.SaveItem(wi)
	}
}
