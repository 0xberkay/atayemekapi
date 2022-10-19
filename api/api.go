package api

import (
	"atayemekapi/database"
	"atayemekapi/models"
	"encoding/json"
	"os"
	"time"

	_ "github.com/goccy/go-json"
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
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
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
		if c.Get("api_key") != database.Api_key {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}
		return c.Next()
	})

	api.Post("/save", Save)

	app.Use(cache.New(cache.Config{
		Expiration:   25 * time.Minute,
		CacheControl: true,
	}))
	api.Get("/announces", GetAllAnnounces)

	api.Get("/add", Save)

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

func Save(c *fiber.Ctx) error {
	p := new(models.AdminData)

	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Bad Request",
			"err":     err,
		})
	}

	if p.Admin != database.Admin {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})

	}

	_, err := database.DB.Collection("foods").UpdateOne(c.Context(), bson.M{"date": p.Date}, bson.M{"$set": bson.M{"menuimage": p.Link}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Saved",
	})
}
