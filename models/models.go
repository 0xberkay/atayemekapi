package models

type Menu struct {
	Link      string
	Date      string
	MenuItems []MenuItem
}

type MenuItem struct {
	Food string
	Gram string
}
