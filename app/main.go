package main

import (
	"os"
	log "github.com/sirupsen/logrus"
)

func main() {
	f := os.Getenv("LUA_FILE")
	if len(f) == 0 {
		f = "aux-addon.lua"
	}

	p := NewParser(f)

	recordedSales, observedBuyouts, err := p.GetData()
	if err != nil {
		log.Fatal(err)
	}

	store := NewPostgres()

	for _, rs := range recordedSales {
		store.SaveRecordedSale(rs)
	}

	for _, ob := range observedBuyouts {
		store.SaveObservedBuyout(ob)
	}
}
