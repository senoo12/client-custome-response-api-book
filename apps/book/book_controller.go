package book

import (
	bookModel "client-response-api-book/infra/db/models"
	infra "client-response-api-book/infra"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
	"fmt"
)

func GetAllBooks(c *fiber.Ctx) error {
	var books []bookModel.Book
	data, err := infra.DB.Query("SELECT * FROM books")

	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	for data.Next(){
		var book bookModel.Book
		if err := data.Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Crated_At, &book.Update_At); err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	if err != nil {
		log.Fatal(err)
	}

	response := infra.ResponseAPIList("Success", "success", fiber.StatusOK, books, len(books))
	return c.JSON(response)
}

func GetBookById(c *fiber.Ctx) error {
	idBooks := c.Params("id")
	var book bookModel.Book

	data, err := infra.DB.Query("SELECT * FROM `books` WHERE id = ?;", idBooks)

	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	if data.Next() {
		if err := data.Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Crated_At, &book.Update_At); err != nil {
			log.Fatal(err)
		}
	} else {
		responseFailed := infra.ResponseAPI("Failed", "failed", fiber.StatusBadRequest, book)
		return c.JSON(responseFailed)
	}

	response := infra.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}

func CreateBook(c *fiber.Ctx) error {
	var book bookModel.Book
	// fmt.Println(book)

	if err := c.BodyParser(&book); err != nil {
		return err
	}

	data, err := infra.DB.Exec("INSERT INTO books(title, author, category) VALUES (?, ?, ?)", book.Title, book.Author, book.Category)
	fmt.Println(book.Title)
	if err != nil {
		log.Fatal(err)
	}

	lastInsertID, err := data.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	book.ID = int(lastInsertID)

	response := infra.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}

func UpdateBook(c *fiber.Ctx) error  {
	idBook := c.Params("id")
	var book bookModel.Book

	if err := c.BodyParser(&book); err != nil {
		return err
	}

	_, err := infra.DB.Exec("UPDATE books SET title = ?, author = ?, category = ?, updated_at = ? WHERE id = ?",
    book.Title, book.Author, book.Category, time.Now(), idBook)
	if err != nil {
		return err
	}

	err = infra.DB.QueryRow("SELECT * FROM books WHERE id = ?;", idBook).
		Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Crated_At, &book.Update_At)

	if err != nil {
		return err
	}

	response := infra.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}

func DeleteBook(c *fiber.Ctx) error  {
	idBook := c.Params("id")
	var book bookModel.Book

	err := infra.DB.QueryRow("SELECT * FROM books WHERE id = ?;", idBook).
	Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Crated_At, &book.Update_At)

	if err != nil {
		responseFailed := infra.ResponseAPI("Failed", "failed", fiber.StatusBadRequest, book)
		return c.JSON(responseFailed)
	}
	_, err = infra.DB.Exec("DELETE FROM books WHERE id = ?", idBook)
	if err != nil {
		return err
	}

	response := infra.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}
