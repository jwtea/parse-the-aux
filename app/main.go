package main

import log "github.com/sirupsen/logrus"

func main() {
	p := NewParser("aux-addon.lua")

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
