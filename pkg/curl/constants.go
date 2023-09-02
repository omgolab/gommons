package gccurl

import "regexp"

// Sample curl header string:
//
//	curl 'https://domain.com/api/charge/instant' \
//	  -X 'POST' \
//	  -H 'accept: application/json, text/javascript, */*; q=0.01' \
//	  -H 'accept-language: en-US,en;q=0.9' \
//	  ........
//	  --compressed
var curlHeaderParser = regexp.MustCompile(`(?m)-H\W+([^:\s]+)\W+([^']+)`)
