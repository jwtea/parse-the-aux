package main

import (
	"bufio"
	"os"
	"strconv"
	"time"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Parser struct {
	scanner *bufio.Scanner
	line    int
}

// NewParser from file
func NewParser(fname string) *Parser {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return &Parser{scanner: scanner, line: 0}
}

// GetItems from the file being scanned
func (p *Parser) GetItems() (items []*DataPoint) {
	if !p.FindHistoryStart() {
		log.Error("Could not find start of history from aux lua")
		return nil
	}

	for p.scanner.Scan() {
		p.line++

		if iRegexp.MatchString(p.scanner.Text()) {
		  items = append(items, parseDataPoints(p.scanner.Text())...)
		}
	}

	return items
}

// FindHistoryStart proceeds through scanner util matching on the history regex
func (p *Parser) FindHistoryStart() bool {
	for p.scanner.Scan() {
		if hRegexp.MatchString(p.scanner.Text()) {
			p.line++
			return true
		}
		p.line++
	}
	return false
}

func defaultReportData() (string, string, string) {
	realm := "mograine"
	faction := "A"
	reporter := "anonymous"
	return realm, faction, reporter
}

func parseBuyoutPoints(line string) ([]*DataPoint, error) {
	realm, faction, reporter := defaultReportData() //todo
	items := make([]*DataPoint, 0)
	matches := buyoutRegexp.FindStringSubmatch(line)
	if len(matches) != 4 || len(matches[3]) == 0 {
		return nil, fmt.Errorf("could not parse buyout point `%s`. parsed: `%+v`", matches[0], matches[1:len(matches)-1])
	}

	itemID := matches[1]

	timestamp, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, err
	}

	buyoutCopper, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, err
	}

	//Todo
	buyoutItem := &DataPoint{
		Time:        time.Unix(timestamp, 0),
		DataType:    "min_buyout",
		Realm:       realm,
		Faction:     faction,
		Reporter:    reporter,
		ItemID:      itemID,
		ItemName:    "Unknown",
		CopperValue: int32(buyoutCopper),
	}

	return append(items, buyoutItem), nil
}

func parseSalePoints(line string) ([]*DataPoint, error) {
	realm, faction, reporter := defaultReportData() // todo
	items := make([]*DataPoint, 0)
	matches := saleRegexp.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("could not parse sale points. parsed: `%+v`", matches)
	}

	itemID := matches[1]
	sales := strings.Split(matches[2], ";")

	for _, s := range sales {
		if len(s) == 0 {
			continue
		}

		splitSale := strings.Split(s, "@")

		if len(splitSale) < 2 {
			return nil, fmt.Errorf("could not split recorded sale `%s`", s)
		}		

		timestamp, err := strconv.ParseInt(splitSale[1], 10, 64)
		if err != nil {
			return nil, err
		}
	
		saleCopper, err := strconv.Atoi(splitSale[0])
		if err != nil {
			return nil, err
		}

		saleItem := &DataPoint{
			Time:        time.Unix(timestamp, 0),
			DataType:    "recorded_sale",
			Realm:       realm,
			Faction:     faction,
			Reporter:    reporter,
			ItemID:      itemID,
			ItemName:    "Unknown",
			CopperValue: int32(saleCopper),
		}

		items = append(items, saleItem)
	}

	return items, nil
}

// parseDataPoints takes a string and checks for regex match on DataPoint history regex
func parseDataPoints(s string) []*DataPoint {
	points, err := parseSalePoints(s)
	if err != nil {
		log.Error(err, s)
	}

	buyoutPoints, err := parseBuyoutPoints(s)
	if err != nil {
		log.Error(err, s)
	}

	for _, b := range buyoutPoints {
		points = append(points, b)
	}

	log.Info(fmt.Sprintf("Parsed %d data points.", len(points)))
	return points
}
