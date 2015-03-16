package main

import (
	"fmt"

	extrict "github.com/odeke-em/extrict/src"
)

func main() {
	matchesChan := extrict.GetAndMatchHttpLinks("https://golang.org")

	for link := range matchesChan {
		fmt.Println(link)
	}
}
