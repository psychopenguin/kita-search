package main

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/psychopenguin/kita-search/pkg/kita"
)

const start string = "https://www.berlin.de/sen/jugend/familie-und-kinder/kindertagesbetreuung/kitas/verzeichnis/ListeKitas.aspx"

func main() {
	db, err := gorm.Open("sqlite3", "kita.db")
	if err != nil {
		panic("Failed to open db")
	}

	defer db.Close()

	db.AutoMigrate(&kita.Kita{})

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
		var k kita.Kita
		k.Permalink = e.Request.URL.String()
		k.Name = e.ChildText("#Allgemein h1")
		k.Email = e.ChildText("#HLinkEMail")
		db.Create(&k)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.URL.String(), r.StatusCode, err)
	})

	c.Visit(start)
	c.Wait()
}
