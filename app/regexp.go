package main

import "regexp"

// itemInfoRegexp matches the beginning of the item information list
var itemInfoRegexp = regexp.MustCompile(`\[\"items\"\]\ \= \{`)

// itemNameRegexp matches the name value from the item information string
var itemNameRegexp = regexp.MustCompile(`"(\w.*)#\|.*\|Hitem:(\d*):`)

// hRegexp matches the beggining of the history
var hRegexp = regexp.MustCompile(`\[\"history\"\]\ \= \{`)

// iRegexp matches the the item entry for history
var iRegexp = regexp.MustCompile(`\[(.*)\] = \"(.*)\"`)

// buyoutRegexp matches the item entry buyout and timestamp for history
var buyoutRegexp = regexp.MustCompile(`\["(\d*):0"\] = "(\d*)#(\d*)#.*\"`)

// saleRegexp matches the sale item entry for history
var saleRegexp = regexp.MustCompile(`\["(\d*):0"\] = "\d*#\d*#(.*)\"`)

// matchCondition used for controlling when scanning should stop
type matchCondition func(s string) bool
