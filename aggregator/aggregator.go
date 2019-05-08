package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var nySections []string
var nytimesXML []byte

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

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	var s SitemapIndex
	var n News

	news_map := make(map[string]NewsMap)
	xml.Unmarshal(nytimesXML, &s) //Unmarshalling the data <created> similar to GET response
	//var itr uint = 0
	for _, Location := range s.Locations {
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)
		for idx, title := range n.Title {
			news_map[title] = NewsMap{n.Description[idx], n.Link[idx]}
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
