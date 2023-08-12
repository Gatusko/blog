package main

import (
	"encoding/xml"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
	"time"
)

type ScrapXml struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
	} `xml:"channel"`
}

func (app *apiConfig) ScrapData(id uuid.UUID, wg *sync.WaitGroup) {
	defer wg.Done()
	feed, err := app.models.Feeds.GetFeedById(id)

	if err != nil {
		log.Println(err)
		return
	}

	response, err := http.Get(feed.Url)
	if err != nil {
		log.Println("response went wrong", err)
		return
	}
	defer response.Body.Close()
	scrppedData := ScrapXml{}
	dec := xml.NewDecoder(response.Body)
	err = dec.Decode(&scrppedData)
	if err != nil {
		log.Printf("Issue decoding : %s", err)
		return
	}
	log.Printf("Succes decoding : %s", scrppedData.Channel.Title)

	err = app.models.Feeds.MarkFeedFetched(id)
	if err != nil {
		log.Println(err)
	}

	log.Println("Scraped succes of Feed", feed)
}

func (apiConfig *apiConfig) workerScrappers(number int, interval time.Duration) {
	ticker := time.NewTicker(interval)
	log.Printf("Enter of workScrapper")
	for range ticker.C {
		feeds, err := apiConfig.models.Feeds.GetNextFeedsToFetch(number)
		if err != nil {
			log.Printf("Error at getting the next Feeds : %s", err)
			continue
		}
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go apiConfig.ScrapData(feed.Id, &wg)
		}
		wg.Wait()
		log.Printf("Succes of scraping the data")
	}
}
