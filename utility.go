package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var getSync = make(chan int, 15)

func getDoc(urlS string) (*goquery.Document, error) {
	getSync <- 1
	defer func() { <-getSync }()

	res, err := http.Get(urlS)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status was %d", res.StatusCode)
	}
	defer res.Body.Close()

	node, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	doc := goquery.NewDocumentFromNode(node)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func tryGetDoc(url string) (*goquery.Document, error) {
	doc, err := getDoc(url)
	count := 0
	for err != nil && count < 3 {
		time.Sleep(3 * time.Second)
		doc, err = getDoc(url)
		count++
	}

	if err != nil {
		log.Error(fmt.Errorf("couldn't load %s: %s", url, err))
	}

	return doc, err
}
