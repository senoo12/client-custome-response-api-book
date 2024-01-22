package user

import (
	common "client-response-api-book/common"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	// "fmt"
	"time"
)

func GetAllUsers(c *fiber.Ctx) error  {
	var users []User
	data, err := common.DB.Query("SELECT * FROM users")

	if err != nil {
		return err
	}

	defer data.Close()

	for data.Next(){
		var user User
		if err := data.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At ); err != nil {
			return err
		}

		users = append(users, user)
	}

	if err != nil {
		return err
	}

	var usersResponse []UserResponse
	for _, user := range users {
		userResponse := UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			Created_At: user.Created_At,
			Updated_At: user.Updated_At,
		}
		usersResponse = append(usersResponse, userResponse)
	}

	response := common.ResponseAPIList("Success", "success", fiber.StatusOK, usersResponse, len(usersResponse))
	return c.JSON(response)
}

func GetUserByID(c *fiber.Ctx) error  {
	idUser := c.Params("id")
	var user User

	data, err := common.DB.Query("SELECT * FROM users WHERE id = ?", idUser)
	if err != nil {
		return err
	}

	defer data.Close()

	if data.Next() {
		if err := data.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At); err != nil {
			return err
		}
	} else {
		responseFailed := common.ResponseAPI("Failed", "failed", fiber.StatusBadRequest, user)
		return c.JSON(responseFailed)
	}

	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.Created_At = user.Created_At
	userResponse.Updated_At = user.Updated_At

	response := common.ResponseAPI("Success", "success", fiber.StatusOK, userResponse)
	return c.JSON(response)
}

func Register(c *fiber.Ctx) error {
	var user User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		responseFailed := common.ResponseAPI("Failed", "error hash", fiber.StatusBadRequest, user)
		return c.JSON(responseFailed)
	}
	user.Password = string(hashedPassword)

	data, err := common.DB.Exec("INSERT INTO users(name, email, password) VALUES(?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	lastInsertID, err := data.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int (lastInsertID)

	response := common.ResponseAPI("Success", "success", fiber.StatusOK, user)
	return c.JSON(response)
}

func Login(c *fiber.Ctx) error  {
	var loginRequest LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to request",
		})
	}

	var user User

	data, err := common.DB.Query("SELECT * FROM `users` WHERE email = ? ", loginRequest.Email)
	if err != nil {
		responseFailed := common.ResponseAPI("Failed", "Password or Email wrong!", fiber.StatusBadRequest, user)
		return c.JSON(responseFailed)
	}

	defer data.Close()

	if data.Next(){
		if err := data.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At); err != nil {
			return err
		}
	} else {
		responseFailed := common.ResponseAPI("Failed", "failed", fiber.StatusBadRequest, user)
		return c.JSON(responseFailed)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		responseFailed := common.ResponseAPI("Failed", "password wrong", fiber.StatusBadRequest, user)
		return c.JSON(responseFailed)
	}

	claims := jwt.MapClaims{
		"email": loginRequest.Email,
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var loginResponse LoginResponse
	loginResponse.ID = user.ID
	loginResponse.Name = user.Name
	loginResponse.Email = user.Email
	loginResponse.Token = t

	return c.JSON(fiber.Map{
		"data": loginResponse,
	})
}	

func UpdateUser(c *fiber.Ctx) error  {
	idUser := c.Params("id")
	var user User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		responseFailed := common.ResponseAPI("Failed", "error hash", fiber.StatusBadRequest, user)
		return c.JSON(responseFailed)
	}
	user.Password = string(hashedPassword)

	_, err = common.DB.Exec("UPDATE users SET name = ?, email = ?, password = ?, updated_at = ? WHERE id = ?", user.Name, user.Email, user.Password, time.Now(), idUser)
	if err != nil {
		return err
	}

	err = common.DB.QueryRow("SELECT * FROM users WHERE id = ?", idUser).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At)
	if err != nil {
		return err
	}

	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.Created_At = user.Created_At
	userResponse.Updated_At = user.Updated_At

	response := common.ResponseAPI("Success", "Success edit", fiber.StatusOK, userResponse)
	return c.JSON(response)
}

func DeleteUser(c *fiber.Ctx) error {
	idUser := c.Params("id")
	var user User

	err := common.DB.QueryRow("SELECT * FROM users WHERE id = ?", idUser).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At)
	if err != nil {
		return err
	}

	_, err = common.DB.Exec("DELETE FROM users WHERE id = ?", idUser)
	if err != nil {
		return err
	}

	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.Created_At = user.Created_At
	userResponse.Updated_At = user.Updated_At

	response := common.ResponseAPI("Success", "success deleted", fiber.StatusOK, userResponse)
	return c.JSON(response)
}