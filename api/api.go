package api

import (
	"atayemekapi/database"
	"atayemekapi/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
)

func ApiRunner() {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		ServerHeader: "Fiber",
		AppName:      "Atayemek API",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	Setup(app)

	//heroku port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}

func Setup(app *fiber.App) {
	api := app.Group("/api", func(c *fiber.Ctx) error {
		if c.Get("api_key") != os.Getenv("api_key") {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}
		return c.Next()
	})

	api.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))
	api.Get("/announces", GetAllAnnounces)
	api.Get("/update", IsUpdateTrue)

	menu := api.Group("/menu")
	menu.Get("/all", GetAllMenu)
	menu.Get("/today", GetTodayMenu)

}

func GetAllAnnounces(c *fiber.Ctx) error {
	announces := []models.Announce{}
	cursor, err := database.DB.Collection("announces").Find(c.Context(), bson.D{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}
	if err = cursor.All(c.Context(), &announces); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{
		"success":   true,
		"announces": announces,
	})
}

func GetAllMenu(c *fiber.Ctx) error {
	menu := []models.Menu{}
	cursor, err := database.DB.Collection("foods").Find(c.Context(), bson.D{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}
	if err = cursor.All(c.Context(), &menu); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"menu":    menu,
	})
}

func GetTodayMenu(c *fiber.Ctx) error {
	menu := models.Menu{}

	todayDataStr := time.Now().Format("02.01.2006")

	err := database.DB.Collection("foods").FindOne(c.Context(), bson.M{"date": todayDataStr}).Decode(&menu)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"menu":    menu,
	})

}

func IsUpdateTrue(c *fiber.Ctx) error {
	update := os.Getenv("update")
	if update == "true" {
		return c.JSON(fiber.Map{
			"success": true,
			"update":  true,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"update":  false,
	})

}
