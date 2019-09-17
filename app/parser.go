package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type WowItem struct {
	Timestamp time.Time
	MinBuyout int32
}

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
func (p *Parser) GetItems() (items []WowItem) {
	if !p.FindHistoryStart() {
		return []WowItem{}
	}

	for p.scanner.Scan() {
		if iRegexp.MatchString(p.scanner.Text()) {
			wi, err := parseToWowItem(p.scanner.Text())
			if err != nil {
				log.Warn("cannot parse item")
			} else {
				items = append(items, *wi)
			}
		}

		p.line++
	}

	return items
}

// FindHistoryStart proceeds through scanner util matching on the history regex
// and returns
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

// parseToWowItem takes a string and checks for regex match on item history regex
func parseToWowItem(s string) (*WowItem, error) {
	m := ibRegexp.FindStringSubmatch(s)
	if len(m) >= 2 {
		ti, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return nil, err
		}

		//todo null values should be stored?
		mb, err := strconv.Atoi(m[2])
		if err != nil {
			mb = 0
		}

		i := WowItem{
			Timestamp: time.Unix(ti, 0),
			MinBuyout: int32(mb),
		}
		return &i, nil
	}

	return nil, errors.New("Invalid data")
}
