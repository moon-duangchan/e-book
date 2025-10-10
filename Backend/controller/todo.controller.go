package controller

import (
    "Backend/database"
    "Backend/models"
    "github.com/gofiber/fiber/v2"
)

//getall
func GetTodos(c *fiber.Ctx) error{
    db := database.DBConn
    var todos []models.Todo
    db.Find(&todos)
    return c.JSON(&todos)
}
//get by id
func GetTodoByID(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var todo models.Todo
	if err := db.First(&todo, id).Error; err != nil {
		return c.Status(404).SendString("Not found")
	}

	return c.JSON(todo)
}


func CreateTodo(c *fiber.Ctx) error {
	db := database.DBConn
	var todo models.Todo

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString("Bad request")
	}

	if err := db.Create(&todo).Error; err != nil {
		return c.Status(500).SendString("Failed to create")
	}

	return c.JSON(todo)
}

func UpdateTodo(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var todo models.Todo
	if err := db.First(&todo, id).Error; err != nil {
		return c.Status(404).SendString("Not found")
	}

	var input models.Todo
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).SendString("Bad request")
	}

	db.Model(&todo).Updates(input)
	return c.JSON(todo)
}


func DeleteTodo(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	if err := db.Delete(&models.Todo{}, id).Error; err != nil {
		return c.Status(500).SendString("Delete failed")
	}

	return c.SendStatus(204)
}
