package main

import "fmt"
import "time"
import "log"
import "strings"
import "github.com/gocolly/colly"

const start string = "https://www.berlin.de/sen/jugend/familie-und-kinder/kindertagesbetreuung/kitas/verzeichnis/ListeKitas.aspx"

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.berlin.de"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*berlin.de*",
		Parallelism: 3,
		Delay:       1 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL.String())
	})

	c.OnHTML("#DataList_Kitas tr a", func(e *colly.HTMLElement) {
		//berlin.de website sends not escaped spaces in it href link
		//this just replace the spaces with + signs to scape them
		l := strings.Replace(e.Request.AbsoluteURL(e.Attr("href")), " ", "+", -1)
		c.Visit(l)
	})
	c.OnHTML("#frmKitaDetailNeu", func(e *colly.HTMLElement) {
		name := e.ChildText("#Allgemein h1")
		fmt.Println(name)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.URL.String(), r.StatusCode, err)
	})

	c.Visit(start)
	c.Wait()
}
