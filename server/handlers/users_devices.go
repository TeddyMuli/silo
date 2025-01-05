package handlers

import (
	"context"
	"log"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
);

//Create a UserDevice
func createUserDevice(c *fiber.Ctx, db *pgxpool.Pool) error {
	var userDevice models.UserDevice

	if err := c.BodyParser(&userDevice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	userDevice.ID = uuid.New().String()
	userDevice.CreatedAt = time.Now()

	query := `
		INSERT INTO userdevices (id, user_organization_id, device_id, created_at)
		VALUES ($1, $2, $3, $4);
	`
	_, err := db.Exec(
		context.Background(),
		query,
		userDevice.ID, userDevice.UserOrganizationID, userDevice.DeviceID, userDevice.CreatedAt,
	)
	
	if err != nil {
		log.Println("Error creating user device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user device", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User Device created successfully!"})
}

//Delete a UserDevice
func deleteUserDevice(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := `
		DELETE FROM userdevices
		WHERE id = $1;
	`

	_, err := db.Exec(context.Background(), query, id)
	if err != nil {
		log.Println("Error deleting User Device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting User Device", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Device deleted successfully!"})
}


//Get a UserDevice
func getUserDevice(c *fiber.Ctx, db *pgxpool.Pool) error {
	user_organization_id := c.Params("user_organization_id");
	device_id := c.Params("device_id");

	if user_organization_id == "" || device_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing parameters"})
	}

	var userDevice models.UserDevice

	query := `
		SELECT id, user_organization_id, device_id, created_at
		FROM devices
		WHERE user_organization_id = $1 AND device_id = $2;
	`
	err := db.QueryRow(
		context.Background(), 
		query, 
		user_organization_id, device_id,
	).Scan(&userDevice.ID, &userDevice.UserOrganizationID, &userDevice.DeviceID, &userDevice.CreatedAt)
	if err != nil {
		log.Println("Error retrieving user device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving user device", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(userDevice)
}

func RegisterUserDeviceRoutes(app *fiber.App, db *pgxpool.Pool) {
	userDeviceGroup := app.Group("/user_device")

	userDeviceGroup.Post("/create", func(c *fiber.Ctx) error {
		return createUserDevice(c, db)
	})
	userDeviceGroup.Get("/fetch/:user_organization_id", func(c *fiber.Ctx) error {
		return getUserDevice(c, db)
	})
	userDeviceGroup.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return deleteUserDevice(c, db)
	})
}
