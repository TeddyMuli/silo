package handlers

import (
	"context"
	"log"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CreateFleet inserts a new fleet into the database
func CreateFleet(c *fiber.Ctx, db *pgxpool.Pool) error {
	var fleet models.Fleet

	if err := c.BodyParser(&fleet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	var existingFleet models.Fleet
	checkQuery := `
		SELECT id FROM fleets
		WHERE name = $1
		AND organization_id = $2
	`

	err := db.QueryRow(
		context.Background(),
		checkQuery,
		fleet.Name,
		fleet.OrganizationID,
	).Scan(&existingFleet.ID)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Fleet with the same name exists!"})
	} else if err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing fleet", "message": err.Error()})
	}

	fleet.ID = uuid.New().String()
	fleet.CreatedAt = time.Now()

	query := `
		INSERT INTO fleets (id, name, organization_id, created_at)
		VALUES ($1, $2, $3, $4);
	`

	_, err = db.Exec(
		context.Background(), 
		query, 
		fleet.ID, fleet.Name, fleet.OrganizationID, fleet.CreatedAt,
	)

	if err != nil {
		log.Println("Error creating fleet:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating fleet", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Fleet created successfully!"})
}

// GetFleetByID retrieves a fleet by its ID
func GetAllFleets(c *fiber.Ctx, db *pgxpool.Pool) error {
	organization_id := c.Params("organization_id")
	var rows pgx.Rows

	if organization_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organization_id missing"})
	}

	query := `
		SELECT id, name, organization_id, created_at
		FROM fleets
		WHERE organization_id = $1;
	`

	rows, err := db.Query(
		context.Background(), 
		query, 
		organization_id,
	)

	if err != nil {
		log.Println("Error retrieving fleets:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving fleets", "message": err.Error()})
	}
	defer rows.Close()

	var fleets []models.Fleet

	for rows.Next() {
		var fleet models.Fleet
		if err := rows.Scan(&fleet.ID, &fleet.Name, &fleet.OrganizationID, &fleet.CreatedAt); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		fleets = append(fleets, fleet)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fleets)
}

// GetFleetByID retrieves a fleet by its ID
func GetFleetByID(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	var fleet models.Fleet

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := `
		SELECT id, name, organization_id, created_at
		FROM fleets
		WHERE id = $1;
	`
	err := db.QueryRow(
		context.Background(), 
		query, 
		id,
	).Scan(&fleet.ID, &fleet.Name, &fleet.OrganizationID, &fleet.CreatedAt)
	if err != nil {
		log.Println("Error retrieving fleet:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving fleet", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fleet)
}

// UpdateFleet updates an existing fleet
func UpdateFleet(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")
	var fleet models.Fleet

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	if err := c.BodyParser(&fleet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	query := `
		UPDATE fleets
		SET name = $2
		WHERE id = $1;
	`
	_, err := db.Exec(
		context.Background(), 
		query, 
		id, fleet.Name,
	)

	if err != nil {
		log.Println("Error updating fleet:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating fleet", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Fleet updated successfully!"})
}

// DeleteFleet deletes a fleet by its ID
func DeleteFleet(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := `
		DELETE FROM fleets
		WHERE id = $1;
	`

	_, err := db.Exec(context.Background(), query, id)

	if err != nil {
		log.Println("Error deleting fleet:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting fleet", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Fleet deleted successfully!"})
}

func RegisterFleetRoutes(app *fiber.App, db *pgxpool.Pool) {
	fleetGroup := app.Group("/fleet")

	fleetGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateFleet(c, db)
	})
	fleetGroup.Put("/update/:id", func(c *fiber.Ctx) error {
		return UpdateFleet(c, db)
	})
	fleetGroup.Get("/fetch/all/:organization_id", func(c *fiber.Ctx) error {
		return GetAllFleets(c, db)
	})
	fleetGroup.Get("/fetch/specific/:id", func(c *fiber.Ctx) error {
		return GetFleetByID(c, db)
	})
	fleetGroup.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return DeleteFleet(c, db)
	})
}
