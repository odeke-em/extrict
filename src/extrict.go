package extrict

import "regexp"

const (
	HttpPattern = `(https?://[^"]+\.[\w\d]+)`
)

func GetAndMatch(uri string, pattern string) chan string {
	return getAndMatch(uri, regexer(pattern))
}

func GetAndMatchHttpLinks(uri string) chan string {
	return GetAndMatch(uri, HttpPattern)
}

func crawl(uri string, pattern *regexp.Regexp, depth uint32, saveChan chan string) {
	// Incomplete
	if depth == 0 {
		return
	}
	if depth >= 1 {
		depth -= 1
	}

	producerChan := getAndMatch(uri, pattern)

	unloadStringChan(producerChan, saveChan)
}
