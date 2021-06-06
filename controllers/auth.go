package controllers

import (
	"os"
	"time"

	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/config"
	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/models"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string, c *fiber.Ctx) (*models.User, error) {
	userCollection := config.MI.DB.Collection(os.Getenv("USER_COLLECTION"))
	query := bson.D{{Key: "email", Value: e}}
	var user models.User
	err := userCollection.FindOne(c.Context(), query).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func Login(c *fiber.Ctx) error {
	// userCollection := config.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	// Get email address

	var input models.User
	var ud models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on login request",
			"error":   err,
		})
	}

	email := input.Email
	password := input.Password
	user, err := getUserByEmail(*email, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on email",
			"error":   err,
		})
	}

	ud = models.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}
	// find User and return

	if !CheckPasswordHash(*password, *ud.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Password",
			"data":    nil,
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

}

func CreateUser(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection(os.Getenv("USER_COLLECTION"))

	user := new(models.User)
	err := c.BodyParser(user)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	user.ID = nil
	hash, err := hashPassword(*user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot hash password",
			"error":   err,
		})
	}

	user.Password = &hash

	result, err := userCollection.InsertOne(c.Context(), user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "cannot insert data",
			"error":   err,
		})
	}

	user = &models.User{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	userCollection.FindOne(c.Context(), query).Decode(user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": user,
		},
	})
}

func ChangePassword(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection(os.Getenv("USER_COLLECTION"))
	var input models.ChangePassword
	var ud models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Passing data error",
			"error":   err,
		})
	}

	email := input.Email
	oldPassword := input.OldPassword
	newPassword := input.NewPassword

	if oldPassword == newPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Old password is same with new password",
		})
	}

	user, err := getUserByEmail(*email, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Can not find user",
			"error":   err,
		})
	}

	ud = models.User{
		Email:    user.Email,
		Password: user.Password,
	}

	if !(CheckPasswordHash(*oldPassword, *ud.Password)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid password",
		})
	}

	newHashPassword, err := hashPassword(*newPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Can not hash password",
			"error":   err,
		})
	}

	ud.Password = &newHashPassword
	update := bson.D{{
		Key: "$set", Value: ud,
	}}

	query := bson.D{{Key: "email", Value: &ud.Email}}
	err = userCollection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
