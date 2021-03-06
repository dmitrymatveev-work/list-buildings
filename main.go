package main

import (
	"fmt"
	"os"
	"sync"

	logpkg "list-buildings/log"
	"list-buildings/model"
	reportpkg "list-buildings/report"
	"list-buildings/web"
)

func main() {
	log := logpkg.New("errors.log")
	report := reportpkg.New("buildings.csv", log)
	webClient := web.NewClient(log)

	streetsDoc, err := webClient.TryGetDoc("https://realt.by/buildings/?utm_source=h-menu&utm_medium=menu&utm_campaign=menu")
	if err != nil {
		os.Exit(1)
	}

	streetURLs := web.ObtainStreetURLs(streetsDoc)

	var wg sync.WaitGroup
	wg.Add(len(streetURLs))

	buildings := make(chan *model.Building, 100)

	for _, streetURL := range streetURLs {
		go func(streetURL string) {
			defer wg.Done()

			streetDoc, err := webClient.TryGetDoc(streetURL)
			if err != nil {
				return
			}

			buildingURLs := web.ObtainBuildingURLs(streetDoc)

			for _, buildingURL := range buildingURLs {
				buildingDoc, err := webClient.TryGetDoc(buildingURL)
				if err != nil {
					continue
				}

				buildings <- web.ObtainBuilding(buildingDoc)
			}
		}(streetURL)
	}

	go func() {
		wg.Wait()
		close(buildings)
	}()

	for b := range buildings {
		if b.IsBrick && b.IsApartment {
			report.Write(b)
			fmt.Println(b)
		}
	}
}
