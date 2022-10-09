package main

import (
	"atayemekapi/database"
	"atayemekapi/models"
	"context"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	database.Connect()
}

func main() {

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".elementor-element-a08ae94", func(e *colly.HTMLElement) {
		e.ForEach("div.post-title > a", func(_ int, e *colly.HTMLElement) {
			e.Text = strings.TrimSpace(e.Text)
			e.Text = strings.Replace(e.Text, " ", "", -1)

			// date, err := time.Parse("02.01.2006", e.Text)
			// if err != nil {
			// 	log.Println(err)
			// }

			//Find date in database
			var menu models.Menu
			err := database.DB.Collection("foods").FindOne(context.Background(), bson.M{"date": e.Text}).Decode(&menu)
			if err != nil {
				log.Println(err)
			}

			if menu.Date == "" {
				link := e.Attr("href")

				menu = models.Menu{
					Link: link,
					Date: e.Text,
				}
				detailCollector := c.Clone()

				detailCollector.OnHTML("tr", func(e *colly.HTMLElement) {
					//if element is not first row
					if e.Index > 0 {
						menuItem := models.MenuItem{}
						count := 0
						e.ForEach("td", func(_ int, e *colly.HTMLElement) {
							e.Text = strings.TrimSpace(e.Text)
							if e.Text != "" {
								if count > 1 {
									return
								} else {
									if count%2 == 0 {
										menuItem.Food = e.Text
									} else if count%2 == 1 {
										menuItem.Gram = e.Text
									}
								}
								count++

							}
						})
						menu.MenuItems = append(menu.MenuItems, menuItem)
					}
				})

				detailCollector.Visit(e.Request.AbsoluteURL(link))
				log.Println(menu)
				_, err = database.DB.Collection("foods").InsertOne(context.TODO(), menu)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("Already in database")
			}

		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit("https://birimler.atauni.edu.tr/saglik-kultur-ve-spor-daire-baskanligi/")

}
