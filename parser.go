package main

import (
	"bufio"
	"os"

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
func (p *Parser) GetItems() (items []Item) {
	if !p.FindHistoryStart() {
		return []Item{}
	}

	for p.scanner.Scan() {
		if iRegexp.MatchString(p.scanner.Text()) {
			log.Infof("Found match")
			m := iRegexp.FindStringSubmatch(p.scanner.Text())
			items = append(items, Item{m[1], m[2]})
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
