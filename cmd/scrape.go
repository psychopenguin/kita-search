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

package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/psychopenguin/kita-search/pkg/kita"
	"github.com/spf13/cobra"
)

const start string = "https://www.berlin.de/sen/jugend/familie-und-kinder/kindertagesbetreuung/kitas/verzeichnis/ListeKitas.aspx"

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		scrape()
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scrape() {
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
