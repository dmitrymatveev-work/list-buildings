package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	logpkg "list-buildings/log"

	"github.com/PuerkitoBio/goquery"
)

var log = logpkg.New("errors.log")

func main() {
	doc, err := tryGetDoc("https://realt.by/buildings/?utm_source=h-menu&utm_medium=menu&utm_campaign=menu")
	if err != nil {
		os.Exit(1)
	}

	buildings := make(chan *Building, 100)

	streets := doc.Find("#sub-menu-content-lists .archive-street-list li a")

	var wg sync.WaitGroup
	length := streets.Length()
	wg.Add(length)

	streets.Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			return
		}
		go processStreet(url, buildings, &wg)
	})

	go func() {
		wg.Wait()
		close(buildings)
	}()

	for b := range buildings {
		writeDataToFile(b)
		fmt.Println(b)
	}
}

func processStreet(url string, buildings chan<- *Building, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := tryGetDoc(url)
	if err != nil {
		return
	}

	doc.
		Find("div.wiki div.wiki-left-item > a:last-of-type").
		Each(func(i int, s *goquery.Selection) {
			url, ok := s.Attr("href")
			if !ok {
				return
			}
			building, err := getBuilding(url)
			if err != nil {
				return
			}
			buildings <- building
		})
}

func getBuilding(url string) (*Building, error) {
	doc, err := tryGetDoc(url)
	if err != nil {
		return nil, err
	}

	props := doc.Find("tr.table-row td.table-row-right")

	materialL := props.FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(strings.ToLower(s.Text()), "кирпич")
	}).
		Length()

	statusL := props.FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(strings.ToLower(s.Text()), "жил")
	}).
		Length()

	if materialL == 0 || statusL == 0 {
		return nil, errors.New("doesn't match")
	}

	building := Building{}
	doc.Find("#c39380 div.bread-content > p > a").
		Each(func(i int, s *goquery.Selection) {
			switch i {
			case 2:
				building.Street = s.Text()
			case 3:
				building.Building = s.Text()
				url, ok := s.Attr("href")
				if !ok {
					return
				}
				building.URL = url
			}
		})

	return &building, nil
}
