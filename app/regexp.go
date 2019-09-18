package main

import "regexp"

var hRegexp = regexp.MustCompile(`\[\"history\"\]\ \= \{`)

//Item Regex conditions
var iRegexp = regexp.MustCompile(`\[(.*)\] = \"(.*)\"`)

// Buyout and timestamp condition
var buyoutRegexp = regexp.MustCompile(`\["(\d*):0"\] = "(\d*)#(\d*)#.*\"`)

// Sale condition
var saleRegexp = regexp.MustCompile(`\["(\d*):0"\] = "\d*#\d*#(.*)\"`)
