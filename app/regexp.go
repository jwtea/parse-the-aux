package main

import "regexp"

var hRegexp = regexp.MustCompile(`\[\"history\"\]\ \= \{`)

//Item Regex conditions
var iRegexp = regexp.MustCompile(`\[(.*)\] = \"(.*)\"`)

// Buyout and timestamp condition
var ibRegexp = regexp.MustCompile(`(\d*)#(\d*)#(\d*)@(\d*);`)
