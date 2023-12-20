package handlers

import (
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/database"
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/models"
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/request"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(u *models.User) (string, string, error) {
	ATclaims := jwt.MapClaims{
		"userId":    u.ID,
		"userName":  u.UserName,
		"userEmail": u.Email,
		"exp":       time.Now().Add(time.Minute * 5).Unix(),
	}

	RTclaims := jwt.MapClaims{
		"userId":    u.ID,
		"userName":  u.UserName,
		"userEmail": u.Email,
		"exp":       time.Now().Add(time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ATclaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RTclaims)

	at, err := accessToken.SignedString([]byte(os.Getenv("AT_SECRET")))
	if err != nil {
		return "", "", err
	}

	rt, err := refreshToken.SignedString([]byte(os.Getenv("RT_SECRET")))
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Data",
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot Hash Password",
		})
	}

	user.Password = string(hashPassword)
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Username Or Email Is Already Exist",
		})
	}

	at, rt, err := CreateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Signed Token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"accessToken":  at,
		"refreshToken": rt,
	})
}

func SignIn(c *fiber.Ctx) error {
	userRequest := new(request.SignIn)
	user := new(models.User)
	if err := c.BodyParser(userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Data",
		})
	}

	database.DB.Where("user_name = ? OR email = ?", userRequest.UserNameOrEmail, userRequest.UserNameOrEmail).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Credential",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Credential",
		})
	}

	at, rt, err := CreateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Signed Token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  at,
		"refreshToken": rt,
	})
}

func ShowProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["userName"].(string)
	id := claims["userId"].(float64)
	email := claims["userEmail"].(string)

	return c.JSON(fiber.Map{
		"id":       id,
		"username": username,
		"email":    email,
	})
}

func Refresh(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["userName"].(string)
	id := claims["userId"].(float64)
	email := claims["userEmail"].(string)

	userData := new(models.User)
	userData.UserName = username
	userData.ID = uint(id)
	userData.Email = email
	at, rt, err := CreateToken(userData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Signed Token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  at,
		"refreshToken": rt,
	})
}
