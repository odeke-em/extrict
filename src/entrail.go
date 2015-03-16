package extrict

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	expb "github.com/odeke-em/exponential-backoff"
)

type consumer func(string) string
type producer func(uri string, pattern *regexp.Regexp) chan string
type regexper func(string) *regexp.Regexp

func regexpGenerator() regexper {
	cache := map[string]*regexp.Regexp{}

	return func(pat string) *regexp.Regexp {
		memoized, ok := cache[pat]
		if ok && memoized != nil {
			return memoized
		}

		// TODO: handle regexp.Compile err
		reg, _ := regexp.Compile(pat)
		cache[pat] = memoized
		return reg
	}
}

var regexer = regexpGenerator()

func responseStringer(result interface{}, err error) (stringified []string) {
	if err != nil {
		return
	}
	resp := result.(*http.Response)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	splits := bytes.Split(body, []byte{'\n'})
	for _, row := range splits {
		stringified = append(stringified, string(row))
	}

	return stringified
}

func getAndMatch(uri string, pattern *regexp.Regexp) chan string {
	// TODO: Read in the value of n
	n := uint32(10)

	// Retry n times
	getter := expb.NewUrlGetter(uri, n)

	matchesChan := make(chan string)

	go func() {
		expb.ExponentialBackOff(getter, func(resp interface{}, err error) {
			lines := responseStringer(resp, err)
			for _, line := range lines {
				if len(line) < 1 {
					continue
				}

				matches := pattern.FindAllString(line, -1)
				for _, match := range matches {
					matchesChan <- match
				}
			}
			close(matchesChan)
		})
	}()

	return matchesChan
}

func unloadStringChan(from chan string, to chan string) {
	for res := range from {
		to <- res
	}
}
