package main

import (
	"log"
	"os"
	"server/database"
	"server/handlers"
	"server/redis_pkg"
	"server/routes"
	"server/spaces"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading env variables")
		}	
	}
	
	// Initialize the database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis_pkg.InitRedis()

	// Initialize digital ocean spaces
	spaces.InitS3()

	// Initialize Fiber router
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:5174,https://alx-silo.vercel.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		MaxAge:           300, // Optional: cache preflight requests for 5 minutes
	}))
	// Register all routes
	routes.RegisterRoutes(app, db, redis_pkg.RedisClient)

	// Set up cron job to delete expired folders
	c := cron.New()
	c.AddFunc("@daily", func() { 
		handlers.DeleteExpiredFolders(db) 
		handlers.DeleteExpiredFiles(db)
	})
	c.Start()
	
	defer c.Stop()
    
	// Start the server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
