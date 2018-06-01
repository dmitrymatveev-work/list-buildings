package web

import (
	"list-buildings/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ObtainURLs method obtains page URLs from a document
func ObtainURLs(doc *goquery.Document, selector string) (urls []string) {
	urls = make([]string, 0)
	doc.
		Find(selector).
		Each(func(i int, s *goquery.Selection) {
			if url, ok := s.Attr("href"); ok {
				urls = append(urls, url)
			}
		})
	return
}

// ObtainBuilding method obtains building details from a document
func ObtainBuilding(doc *goquery.Document) *model.Building {
	building := model.Building{}

	props := doc.Find("tr.table-row td.table-row-right")

	building.IsBrick = props.FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(strings.ToLower(s.Text()), "кирпич")
	}).
		Length() > 0

	building.IsApartment = props.FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(strings.ToLower(s.Text()), "жил")
	}).
		Length() > 0

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

	return &building
}
