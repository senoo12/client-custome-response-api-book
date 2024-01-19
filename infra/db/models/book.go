package models

import (
	"time"
)

type Book struct {
	ID int 
	Title string
	Author string
	Category string
	Crated_At time.Time
	Update_At time.Time
}