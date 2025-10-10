package main

import (
    "Backend/models"
    "Backend/database"
    "Backend/controller"
    "Backend/auth"
    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "os"
    "fmt"
)

func initDatabase() {
    var err error

    host := "localhost"
    user := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    dbname := os.Getenv("POSTGRES_DB")
    port := os.Getenv("HOST_DB_PORT")

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        host, user, password, dbname, port,
    )

    database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(" Failed to connect to database:", err)
    }

    fmt.Println(" Database connected successfully!")

    err = database.DBConn.AutoMigrate(&models.Todo{}, &models.User{})
    if err != nil {
        log.Fatal(" Failed to migrate database:", err)
    }

    fmt.Println(" Database migrated successfully!")
}

func setupRoutes(app *fiber.App) {
    // Protect all /todos* routes with JWT auth
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



func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    initDatabase()

    app := fiber.New()
    setupRoutes(app)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋! DB: " + os.Getenv("POSTGRES_DB"))
    })

    if err := app.Listen(":" + port); err != nil {
        log.Fatal(err)
    }
}
