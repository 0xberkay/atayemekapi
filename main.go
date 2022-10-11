package main

import (
	"atayemekapi/api"
	"atayemekapi/database"
	"atayemekapi/helper"
)

func init() {
	database.Connect()
}

func main() {
	go helper.TickerForScraping()
	api.ApiRunner()

}
