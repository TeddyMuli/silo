package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/models"
	"server/redis_pkg"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DeviceData struct {
	SerialNumber string `json:"serial_number"`
	Battery      string `json:"battery"`
	Online       bool   `json:"online"`
}


// CreateDevice inserts a new device into the database
func CreateDevice(c *fiber.Ctx, db *pgxpool.Pool) error {
	var device models.Device

	if err := c.BodyParser(&device); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	if device.Name == "" || device.SerialNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields!"})
	}

	if *device.FleetID == "" {
		device.FleetID = nil
	}

	var existingDevice models.Device
	checkQuery := `
		SELECT id from devices
		WHERE LOWER(name) = LOWER($1)
		AND organization_id = $2
	`

	err := db.QueryRow(
		context.Background(),
		checkQuery,
		device.Name,
		device.OrganizationID,
	).Scan(&existingDevice.ID)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Device with the same name exists!"})
	} else if err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing device", "message": err.Error()})
	}

	checkQuery = `
		SELECT id from devices
		WHERE serial_number = $1
	`

	err = db.QueryRow(
		context.Background(),
		checkQuery,
		device.SerialNumber,
	).Scan(&existingDevice.ID)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Device with the same serial number exists!"})
	} else if err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing device", "message": err.Error()})
	}
	
	device.ID = uuid.New().String()
	device.CreatedAt = time.Now()

	query := `
		INSERT INTO devices (id, name, serial_number, fleet_id, ip_address, created_at, organization_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err = db.Exec(
		context.Background(), 
		query, 
		device.ID, device.Name, device.SerialNumber, device.FleetID, device.IPAddress, device.CreatedAt, device.OrganizationID,
	)

	if err != nil {
		log.Println("Error creating device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating device", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Device created successfully!"})
}

// GetDeviceByID retrieves a device by its ID
func GetDeviceByID(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	var device models.Device

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := `
		SELECT id, name, serial_number, fleet_id, ip_address, created_at
		FROM devices
		WHERE id = $1;
	`
	err := db.QueryRow(
		context.Background(), 
		query, 
		id,
	).Scan(&device.ID, &device.Name, &device.SerialNumber, &device.FleetID, &device.IPAddress, &device.CreatedAt)

	if err != nil {
		log.Println("Error retrieving device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving device", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(device)
}

// GetAllDevices retrieves a device by its organization_id
func GetAllDevices(c *fiber.Ctx, db *pgxpool.Pool) error {
	organization_id := c.Params("organization_id")
	fleet_id := c.Params("fleet_id")

	var query string
	var err error
	var rows pgx.Rows

	if organization_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organization_id missing"})
	}

	if organization_id != "" && fleet_id == "" {
		query = `
			SELECT id, name, serial_number, fleet_id, ip_address, created_at
			FROM devices
			WHERE organization_id = $1;
		`
		rows, err = db.Query(
			context.Background(), 
			query, 
			organization_id,
		)
	} else if organization_id != "" && fleet_id != "" {
		query = `
			SELECT id, name, serial_number, fleet_id, ip_address, created_at
			FROM devices
			WHERE fleet_id = $1;
		`
		rows, err = db.Query(
			context.Background(), 
			query, 
			fleet_id,
		)
	}

	if err != nil {
		log.Println("Error retrieving devices:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving devices", "message": err.Error()})
	}
	defer rows.Close()

	var devices []models.Device

	for rows.Next() {
		var device models.Device
		if err := rows.Scan(&device.ID, &device.Name, &device.SerialNumber, &device.FleetID, &device.IPAddress, &device.CreatedAt); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(devices)
}

// UpdateDevice updates an existing device
func UpdateDevice(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	var device models.Device

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	if err := c.BodyParser(&device); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	var updateFields []string
	var args []interface{}
	argIndex := 1

	if device.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, device.Name)
		argIndex++
	}

	if device.SerialNumber != "" {
		updateFields = append(updateFields, fmt.Sprintf("serial_number = $%d", argIndex))
		args = append(args, device.SerialNumber)
		argIndex++
	}

	if device.FleetID != nil && *device.FleetID != "" {
		updateFields = append(updateFields, fmt.Sprintf("fleet_id = $%d", argIndex))
		args = append(args, device.FleetID)
		argIndex++
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE devices
		SET %s
		WHERE id = $%d;
	`, strings.Join(updateFields, ", "), argIndex)

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error beginning transaction", "message": err.Error()})
	}

	defer tx.Rollback(context.Background())

	_, err = db.Exec(
		context.Background(), 
		query, 
		args...,
	)

	if err != nil {
		log.Println("Error updating device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating device", "message": err.Error()})
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Println("Error committing transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error committing transaction", "message": err.Error()})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Device updated successfully!"})
}

// DeleteDevice deletes a device by its ID
func DeleteDevice(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := `
		DELETE FROM devices
		WHERE id = $1;
	`
	_, err := db.Exec(context.Background(), query, id)
	if err != nil {
		log.Println("Error deleting device:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting device", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Device deleted successfully!"})
}

func CheckDevice(c *fiber.Ctx, db *pgxpool.Pool, redisClient *redis.Client) error {
	var deviceData DeviceData
	if err := c.BodyParser(&deviceData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	if deviceData.SerialNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing serial number"})
	}

	ctx := context.Background()
	deviceDataJSON, err := json.Marshal(deviceData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error marshaling device data", "message": err.Error()})
	}

	err = redis_pkg.RedisClient.Set(ctx, fmt.Sprintf("device:%s", deviceData.SerialNumber), deviceDataJSON, time.Minute*3).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving device data to Redis", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Device data saved successfully"})
}

// GetDeviceData retrieves the device data from Redis
func GetDeviceData(c *fiber.Ctx, redisClient *redis.Client) error {
	serialNumber := c.Params("serial_number")
	if serialNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing serial number"})
	}

	ctx := context.Background()
	data, err := redisClient.Get(ctx, fmt.Sprintf("device:%s", serialNumber)).Result()
	if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Device not found"})
	} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving device data from Redis", "message": err.Error()})
	}

	var deviceData DeviceData
	if err := json.Unmarshal([]byte(data), &deviceData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error unmarshaling device data", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(deviceData)
}

func RegisterDeviceRoutes(app *fiber.App, db *pgxpool.Pool, redisClient *redis.Client) {
	deviceGroup := app.Group("/device")

	deviceGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateDevice(c, db)
	})
	deviceGroup.Put("/update/:id", func(c *fiber.Ctx) error {
		return UpdateDevice(c, db)
	})
	deviceGroup.Get("/fetch/all/:organization_id/:fleet_id?", func(c *fiber.Ctx) error {
		return GetAllDevices(c, db)
	})
	deviceGroup.Get("/fetch/specific/:id", func(c *fiber.Ctx) error {
		return GetDeviceByID(c, db)
	})
	deviceGroup.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return DeleteDevice(c, db)
	})
	deviceGroup.Post("/check-in", func(c *fiber.Ctx) error {
		return CheckDevice(c, db, redisClient)
	})
	deviceGroup.Get("/fetch/data/:serial_number", func(c *fiber.Ctx) error {
		return GetDeviceData(c, redisClient)
	})
}
