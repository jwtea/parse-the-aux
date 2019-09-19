package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	errors "golang.org/x/xerrors"

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

// GetData returns buyout items and recorded sales
func (p *Parser) GetData() (recordedSales []*RecordedSale, observedBuyouts []*ObservedBuyout, err error) {
	if !p.FindHistoryStart() {
		return nil, nil, errors.New("Could not find start of history from aux lua")
	}

	for p.scanner.Scan() {
		p.line++

		if iRegexp.MatchString(p.scanner.Text()) {
			rs, pErr := parseRecordedSalePoints(p.scanner.Text())
			if pErr != nil {
				log.Warn(pErr)
			} else {
				recordedSales = append(recordedSales, rs...)
			}

			bo, pErr := parseBuyoutPoints(p.scanner.Text())
			if pErr != nil {
				log.Warn(pErr)
			} else {
				observedBuyouts = append(observedBuyouts, bo...)
			}
		}
	}

	return recordedSales, observedBuyouts, err
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

func parseBuyoutPoints(line string) ([]*ObservedBuyout, error) {
	items := make([]*ObservedBuyout, 0)
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

	buyoutItem := &ObservedBuyout{
		TimescaleRecord{
			Time:        time.Unix(timestamp, 0),
			ItemID:      itemID,
			ItemName:    "Unknown",
			CopperValue: int32(buyoutCopper),
		},
	}

	return append(items, buyoutItem), nil
}

func parseRecordedSalePoints(line string) ([]*RecordedSale, error) {
	items := make([]*RecordedSale, 0)
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

		saleItem := &RecordedSale{
			TimescaleRecord{
				Time:        time.Unix(timestamp, 0),
				ItemID:      itemID,
				ItemName:    "Unknown",
				CopperValue: int32(saleCopper),
			},
		}

		items = append(items, saleItem)
	}

	return items, nil
}
