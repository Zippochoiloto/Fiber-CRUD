package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/config"
	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/models"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	// Get email address

	user := new(models.User)
	c.BodyParser(user)
	// password := c.BodyParser(user.Password)

	fmt.Print(user)

	// find User and return
	query := bson.D{{Key: "Email", Value: user.Email}}

	err := userCollection.FindOne(c.Context(), query).Decode(user)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Can not find user",
			"error":   err,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = user.Email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config(os.Getenv("JWT_KEY"))))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":  user,
			"token": t,
		},
	})

	// func CreateUser (c *fiber.Ctx) error {

	// }
}
