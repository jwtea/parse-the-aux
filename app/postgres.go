package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type TimescaleRecord struct {
	ID          uint32 `gorm:"primary_key"`
	Time        time.Time
	ItemID      string
	ItemName    string
	CopperValue int32 // 1 = 1 copper
}

type RecordedSale struct {
	TimescaleRecord
}

type ObservedBuyout struct {
	TimescaleRecord
}

type DataPoint struct {
	Time     time.Time
	Realm    string // realm data was recorded on
	Faction  string // auction house faction: [A, H, N]
	Reporter string // uploader user or character, havent decided
}

type Postgres struct {
	Conn *gorm.DB
}

//TODO
func defaultReportData() (string, string, string) {
	realm := "mograine"
	faction := "A"
	reporter := "anonymous"
	return realm, faction, reporter
}

// NewPostgres returns a postgres client with default connection string
func NewPostgres() *Postgres {
	gorm, err := gorm.Open("postgres", "postgres://timescale:password@localhost/timescale?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	p := &Postgres{
		Conn: gorm,
	}

	p.migrate()

	return p
}

// migrate db models
func (p *Postgres) migrate() {
	if err := p.Conn.AutoMigrate(&RecordedSale{}).Error; err != nil {
		log.Fatalf("Cannot migrate recorded sale: %s", err)
	}

	if err := p.Conn.AutoMigrate(&ObservedBuyout{}).Error; err != nil {
		log.Fatalf("Cannot migrate observed buyout: %s", err)
	}
}

// SaveRecordedSale for a given timescale record to pg
func (p *Postgres) SaveRecordedSale(r *RecordedSale) {
	p.saveItem(r)
}

// SaveObservedBuyout for a given timescale record to pg
func (p *Postgres) SaveObservedBuyout(r *ObservedBuyout) {
	p.saveItem(r)
}

func (p *Postgres) saveItem(r interface{}) {
	if err := p.Conn.Create(r).Error; err != nil {
		log.Error(err)
	}
}
