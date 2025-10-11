
package route

import (
    "Backend/controller"
    "Backend/auth"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // auth protec todo
    app.Use("/todos", auth.RequireAuth())

    // Todo routes
    app.Get("/todos", controller.GetTodos)
    app.Get("/todos/:id", controller.GetTodoByID)
    app.Post("/todos", controller.CreateTodo)
    app.Put("/todos/:id", controller.UpdateTodo)
    app.Delete("/todos/:id", controller.DeleteTodo)

    // Auth routes
    app.Post("/auth/register", controller.Register)
    app.Get("/auth/verify", controller.VerifyEmail)
    app.Post("/auth/login", controller.Login)
}
