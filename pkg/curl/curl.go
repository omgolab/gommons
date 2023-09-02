package gccurl

import (
	"net/http"
)

// CurlReqToHttpHeaders converts a full curl request string to a map of http headers
// Sample curl header string:
//
//	curl 'https://domain.com/api/charge/instant' \
//	  -X 'POST' \
//	  -H 'accept: application/json, text/javascript, */*; q=0.01' \
//	  -H 'accept-language: en-US,en;q=0.9' \
//	  ........
//	  --compressed
func CurlReqToHttpHeaders(curlHeader string) http.Header {
	headers := make(http.Header)
	for _, match := range curlHeaderParser.FindAllStringSubmatch(curlHeader, -1) {
		headers.Set(match[1], match[2])
	}
	return headers
}
