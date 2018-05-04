package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var getSync = make(chan int, 15)

var errorsC = make(chan error, 100)

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
		errorsC <- fmt.Errorf("couldn't load %s: %s", url, err)
	}

	return doc, err
}

func writeDataToFile(b *Building) {
	f, err := os.OpenFile("buildings.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s,%s,%s\n", b.Street, b.Building, b.URL))
	f.Sync()
}

func writeErrorToFile(e error) {
	f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s\n", e.Error()))
	f.Sync()
}
