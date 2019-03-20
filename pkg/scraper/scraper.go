// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scraper

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/psychopenguin/kita-search/pkg/kita"
)

const (
	CSSPathKitaList    = "#DataList_Kitas tr a"
	CSSPathKitaDetails = "#frmKitaDetailNeu"
)

func Scrape(url string) {

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		Delay:       1 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL.String())
	})

	c.OnHTML(CSSPathKitaList, func(e *colly.HTMLElement) {
		//berlin.de website sends not escaped spaces in it href link
		//this just replace the spaces with + signs to scape them
		l := strings.Replace(e.Request.AbsoluteURL(e.Attr("href")), " ", "+", -1)
		c.Visit(l)
	})
	c.OnHTML(CSSPathKitaDetails, func(e *colly.HTMLElement) {
		var k kita.Kita
		k.Permalink = e.Request.URL.String()
		k.Name = e.ChildText("#Allgemein h1")
		k.Email = e.ChildText("#HLinkEMail")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.URL.String(), r.StatusCode, err)
	})

	c.Visit(url)
	c.Wait()
}
