package main

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
	idb "github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	c idb.Client
	o StoreOptions
}

//StoreOptions config
type StoreOptions struct {
	Address string
	DBName  string
}

// NewStoreOptions with defaults
func NewStoreOptions() StoreOptions {
	return StoreOptions{
		Address: "localhost:8086",
		DBName:  "aux",
	}
}

// NewStore configured from store option values
func NewStore(o StoreOptions) *Store {
	c, err := idb.NewHTTPClient(idb.HTTPConfig{
		Addr: o.Address,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Store{c, o}
}

// SaveItem for a given wow item save to influx
func (s *Store) SaveItem(i WowItem) {
	bp, err := idb.NewBatchPoints(client.BatchPointsConfig{
		Database: s.o.DBName,
	})

	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{}
	var fields map[string]interface{}

	pt, err := idb.NewPoint("items", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoint(pt)

	if err := s.c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
