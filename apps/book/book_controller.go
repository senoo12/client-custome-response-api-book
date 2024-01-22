package book

import (
	common "client-response-api-book/common"
	"github.com/gofiber/fiber/v2"
	"time"
)

func GetAllBooks(c *fiber.Ctx) error {
	var books []Book
	data, err := common.DB.Query("SELECT * FROM books")

	if err != nil {
		return err
	}

	defer data.Close()

	for data.Next(){
		var book Book
		if err := data.Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Created_At, &book.Update_At); err != nil {
			return err
		}
		books = append(books, book)
	}

	if err != nil {
		return err
	}

	response := common.ResponseAPIList("Success", "success", fiber.StatusOK, books, len(books))
	return c.JSON(response)
}

func GetBookById(c *fiber.Ctx) error {
	idBooks := c.Params("id")
	var book Book

	data, err := common.DB.Query("SELECT * FROM `books` WHERE id = ?;", idBooks)

	if err != nil {
		return err
	}
	defer data.Close()

	if data.Next() {
		if err := data.Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Created_At, &book.Update_At); err != nil {
			return err
		}
	} else {
		responseFailed := common.ResponseAPI("Failed", "failed", fiber.StatusBadRequest, book)
		return c.JSON(responseFailed)
	}

	response := common.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}

func CreateBook(c *fiber.Ctx) error {
	var book Book

	if err := c.BodyParser(&book); err != nil {
		return err
	}

	data, err := common.DB.Exec("INSERT INTO books(title, author, category) VALUES (?, ?, ?)", book.Title, book.Author, book.Category)
	if err != nil {
		return err
	}

	lastInsertID, err := data.LastInsertId()
	if err != nil {
		return err
	}
	book.ID = int(lastInsertID)

	response := common.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}

func UpdateBook(c *fiber.Ctx) error  {
	idBook := c.Params("id")
	var book Book

	if err := c.BodyParser(&book); err != nil {
		return err
	}

	_, err := common.DB.Exec("UPDATE books SET title = ?, author = ?, category = ?, updated_at = ? WHERE id = ?",
    book.Title, book.Author, book.Category, time.Now(), idBook)
	if err != nil {
		return err
	}

	err = common.DB.QueryRow("SELECT * FROM books WHERE id = ?;", idBook).
		Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Created_At, &book.Update_At)

	if err != nil {
		return err
	}

	response := common.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}

func DeleteBook(c *fiber.Ctx) error  {
	idBook := c.Params("id")
	var book Book

	err := common.DB.QueryRow("SELECT * FROM books WHERE id = ?;", idBook).
	Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Created_At, &book.Update_At)

	if err != nil {
		responseFailed := common.ResponseAPI("Failed", "failed", fiber.StatusBadRequest, book)
		return c.JSON(responseFailed)
	}
	_, err = common.DB.Exec("DELETE FROM books WHERE id = ?", idBook)
	if err != nil {
		return err
	}

	response := common.ResponseAPI("Success", "success", fiber.StatusOK, book)
	return c.JSON(response)
}
