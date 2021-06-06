package routes

import (
	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/controllers"
	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/middleware"
	"github.com/gofiber/fiber/v2"
)

func TodoRoute(route fiber.Router) {
	route.Get("", middleware.Protected(), controllers.GetTodos)
	route.Post("", middleware.Protected(), controllers.CreateTodos)
	route.Put("/:id", middleware.Protected(), controllers.UpdateTodo)
	route.Delete("/:id", middleware.Protected(), controllers.DeleteTodo)
	route.Get("/:id", middleware.Protected(), controllers.GetTodo)
}

func UserRoute(route fiber.Router) {
	route.Post("login", controllers.Login)
	route.Post("register", controllers.CreateUser)
	route.Post("change-password", controllers.ChangePassword)
}
