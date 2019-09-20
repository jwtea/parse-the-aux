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
	idNameMap, ok := p.FindIDNameMap(hRegexp.MatchString)
	if !ok {
		return nil, nil, errors.New("Cannot find start of items from aux lua")
	}

	for p.scanner.Scan() {
		p.line++
		if iRegexp.MatchString(p.scanner.Text()) {
			rs, pErr := parseRecordedSalePoints(idNameMap, p.scanner.Text())
			if pErr != nil {
				log.Warn(pErr)
			} else {
				recordedSales = append(recordedSales, rs...)
			}

			bo, pErr := parseBuyoutPoints(idNameMap, p.scanner.Text())
			if pErr != nil {
				log.Warn(pErr)
			} else {
				observedBuyouts = append(observedBuyouts, bo)
			}
		}
	}

	return recordedSales, observedBuyouts, err
}

// FindIDNameMap returns a map of any name -> id data found in file until
// the stop condition.
func (p *Parser) FindIDNameMap(stopCheck matchCondition) (idNameMap map[string]string, ok bool) {
	idNameMap = make(map[string]string)
	for p.scanner.Scan() {
		p.line++
		if itemNameRegexp.MatchString(p.scanner.Text()) {
			id, name, pErr := parseItemName(p.scanner.Text())
			if pErr != nil {
				log.Warn(pErr)
			} else {
				idNameMap[id] = name
			}
		}
		// Scan until we find our match condition
		if stopCheck(p.scanner.Text()) == true {
			break
		}
	}

	if len(idNameMap) == 0 {
		return idNameMap, false
	}

	return idNameMap, true
}

func parseBuyoutPoints(idNameMap map[string]string, line string) (*ObservedBuyout, error) {
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
			ID:          getHash32(strings.Join([]string{matches[2], itemID, matches[3]}, "-")),
			Time:        time.Unix(timestamp, 0),
			ItemID:      itemID,
			ItemName:    idNameMap[itemID],
			CopperValue: int32(buyoutCopper),
		},
	}

	return buyoutItem, nil
}

func parseRecordedSalePoints(idNameMap map[string]string, line string) ([]*RecordedSale, error) {
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
				ID:          getHash32(strings.Join([]string{splitSale[1], itemID, splitSale[0]}, "-")),
				Time:        time.Unix(timestamp, 0),
				ItemID:      itemID,
				ItemName:    idNameMap[itemID],
				CopperValue: int32(saleCopper),
			},
		}

		items = append(items, saleItem)
	}

	return items, nil
}

func parseItemName(line string) (string, string, error) {
	matches := itemNameRegexp.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("could not parse item name. parsed: `%+v`", matches)
	}
	return matches[2], matches[1], nil
}
