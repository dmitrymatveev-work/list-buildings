package web

import (
	"fmt"
	"net/http"
	"time"

	"list-buildings/log"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Client is the type that implements specific web-client functionality
type Client struct {
	log   *log.Log
	queue chan int
}

// NewClient method that constructs a new instance of the web-client
func NewClient(log *log.Log) *Client {
	return &Client{log, make(chan int, 15)}
}

func (c *Client) getDoc(url string) (*goquery.Document, error) {
	c.queue <- 1
	defer func() { <-c.queue }()

	res, err := http.Get(url)
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

// TryGetDoc method tries to get a document on a specified resource address
func (c *Client) TryGetDoc(url string) (*goquery.Document, error) {
	doc, err := c.getDoc(url)
	count := 0
	for err != nil && count < 3 {
		time.Sleep(3 * time.Second)
		doc, err = c.getDoc(url)
		count++
	}

	if err != nil {
		c.log.Error(fmt.Errorf("couldn't load %s: %s", url, err))
	}

	return doc, err
}
