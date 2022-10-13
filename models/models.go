package models

type Menu struct {
	Link      string
	Date      string
	MenuItems []MenuItem
	TotelGram int
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
