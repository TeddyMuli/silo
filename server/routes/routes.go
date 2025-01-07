package routes

import (
	"server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, redisClient *redis.Client) {
	// Auth routes
	handlers.RegisterAuthRoutes(app, db)
	// Folder routes
	handlers.RegisterFolderRoutes(app, db)
	// File routes
	handlers.RegisterFileRoutes(app, db)
	// Organization routes
	handlers.RegisterOrganizationRoutes(app, db)
	// Device routes
	handlers.RegisterDeviceRoutes(app, db, redisClient)
	// Fleet routes
	handlers.RegisterFleetRoutes(app, db)
	// Plan routes
	handlers.RegisterPlanRoutes(app, db)
	// User Organization routes
	handlers.RegisterUserOrganizationRoutes(app, db)
	// User Device Routes
	handlers.RegisterUserDeviceRoutes(app, db)
	// Subscription routes
	handlers.RegisterBinRoutes(app, db)
}
