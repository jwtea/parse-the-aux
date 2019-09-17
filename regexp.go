package main

import "regexp"

var hRegexp = regexp.MustCompile(`\[\"history\"\]\ \= \{`)
var iRegexp = regexp.MustCompile(`\[(.*)\] = \"(.*)\"`)
