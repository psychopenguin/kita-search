package main

import "fmt"
import "time"
import "log"
import "net/url"
import "github.com/gocolly/colly"

const start string = "https://www.berlin.de/sen/jugend/familie-und-kinder/kindertagesbetreuung/kitas/verzeichnis/ListeKitas.aspx"

func main(){
	c := colly.NewCollector(
		colly.AllowedDomains("www.berlin.de"),
		colly.Async(true),
	)
	
	c.Limit(&colly.LimitRule{
		DomainGlob: "*berlin.de*",
		Parallelism: 3,
		Delay: 1 * time.Second,
		RandomDelay: 2 * time.Second,
	})
	
	k := c.Clone()

	c.OnRequest(func (r *colly.Request){
		log.Println("Retrieving Kita list from:", r.URL.String())
	})

	c.OnHTML("#DataList_Kitas tr a", func(e *colly.HTMLElement) {
		// URLs returned by this contains spaces and need to be encoded properly
		raw_link := e.Request.AbsoluteURL(e.Attr("href"))
		u, _ := url.Parse(raw_link)
		u.RawQuery = url.QueryEscape(u.RawQuery)
		fmt.Println(u.String())
		k.Visit(u.String())
		k.Wait()
	})
	
	
	k.OnHTML("#frmKitaDetailNeu", func(e *colly.HTMLElement) {
		name := e.ChildText("#Allgemein")
		fmt.Println(name)
	})

	k.OnError(func(r *colly.Response, err error){
		log.Println(r.StatusCode, err)
	})

	
	c.Visit(start)
	c.Wait()
}
