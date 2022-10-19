package models

type Menu struct {
	Link      string
	Date      string
	MenuItems []MenuItem
	TotelGram int
	MenuImage string
}

type Announce struct {
	Link  string
	Title string
	Date  string
	Text  string
}

type MenuItem struct {
	Food string
	Gram string
}

type AdminData struct {
	Admin string `json:"admin"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}
