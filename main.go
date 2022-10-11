package main

import (
	"atayemekapi/api"
	"atayemekapi/database"
)

func init() {
	database.Connect()
}

func main() {
	// helper.Scrapper()
	api.ApiRunner()

}
