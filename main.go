package main

import (
	"fmt"
	"os"
	"strconv"

	extrict "github.com/odeke-em/extrict/src"
)

func main() {
	argc := len(os.Args)
	if argc <= 1 {
		fmt.Fprintf(os.Stderr, "expecting <url> <ext>")
		os.Exit(-1)
	}

	ext := "mp4"
	url := os.Args[1]
	depth := int32(4) // Arbitrary value

	if argc >= 3 {
		ext = os.Args[2]
		if argc >= 4 {
			d, err := strconv.ParseInt(os.Args[3], 10, 32)
			if err == nil {
				depth = int32(d)
			}
		}
	}

	// fmt.Println("crawling", url, ext, depth)
	matches := extrict.CrawlAndMatchByExtension(url, ext, depth)

	for match := range matches {
		fmt.Println(match)
	}
}
