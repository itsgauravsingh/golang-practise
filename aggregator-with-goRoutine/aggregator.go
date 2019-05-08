/*
	This program fetches sitemap from nytimes And Futher looks for latest news from various sections
*/

package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var nySections []string
var nytimesXML []byte
var wg sync.WaitGroup
var queue chan News

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Title         []string `xml:"channel>item>title"`
	Description   []string `xml:"channel>item>description"`
	Link          []string `xml:"channel>item>link"`
	Author        []string `xml:"channel>item>author"`
	PublishedDate []string `xml:"channel>item>pubDate"`
}

type NewsMap struct {
	//Title string
	Description string
	Link        string
}

type NewsAggPage struct {
	Title string
	News  map[string]NewsMap
}

func ini(start bool) bool {
	nySections = []string{"world", "us", "politics", "nyregion", "business", "technology", "science", "health", "sports", "education", "obituaries"}
	nytimesXML = []byte(`<sitemapindex>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/world/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/us/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/politics/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/nyregion/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/business/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/technology/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/sports/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/science/rss.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://www.nytimes.com/svc/collections/v1/publish/https://www.nytimes.com/section/education/rss.xml</loc>
	</sitemap>
</sitemapindex>`)
	log.Println("Inside ini() method")
	return true
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> Hello world </h1>")
}

func newsRoutine(q chan News, Location string) {
	defer wg.Done()
	var n News
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()
	q <- n
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	var s SitemapIndex

	news_map := make(map[string]NewsMap)
	xml.Unmarshal(nytimesXML, &s) //Unmarshalling the data <created> similar to GET response

	queue := make(chan News, 30) // Creating a channel buffer of length 30 to hold the outputs from different goRoutines

	for _, Location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, Location)
	}
	wg.Wait()    // Waiting for goRoutine completion
	close(queue) //closing the channel once the data production has been completed.

	for elem := range queue { // Each element of queue is output w.r.t. one goRoutine execution i.e. Data from each sitemap
		for idx, title := range elem.Title { // Parsing the data received from one of the sitemap
			news_map[title] = NewsMap{elem.Description[idx], elem.Link[idx]}
		}
	}

	p := NewsAggPage{Title: "Amazing News aggregator", News: news_map}
	t, _ := template.ParseFiles("aggregator-template.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Rendered an HTML page")
}

func main() {
	initialised := ini(true)
	if !initialised {
		fmt.Println("Initilise Unsuccessfull")
	}
	fmt.Println("Initilise Successfully")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":9090", nil)
}
