package book

import (
	"time"
)

type Book struct {
	ID int 
	Title string
	Author string
	Category string
	Created_At time.Time
	Update_At time.Time
}