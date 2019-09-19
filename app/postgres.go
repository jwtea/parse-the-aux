package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DataPoint struct {
	Time        time.Time
	Realm       string // realm data was recorded on
	Faction     string // auction house faction: [A, H, N]
	Reporter    string // uploader user or character, havent decided
	DataType    string // [observed_buyout, reported_sale]
	ItemID      string
	ItemName    string
	CopperValue int32 // 1 = 1 copper
}

type Postgres struct {
	Conn *gorm.DB
}

// NewStore configured from store option values
func NewPostgres() *Postgres {
	gorm, err := gorm.Open("postgres", "postgres://timescale:password@localhost/timescale?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	gorm.AutoMigrate(&DataPoint{})

	return &Postgres{
		Conn: gorm,
	}
}

// SaveItem for a given wow item save to influx
func (p *Postgres) SaveItem(i *DataPoint) {
	p.Conn.Create(i)
}
