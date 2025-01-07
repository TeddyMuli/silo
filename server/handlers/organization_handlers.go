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

//Creates an organization
func CreateOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	// Define a struct to capture both Organization and User ID from the request body
	type requestData struct {
		Name   string `json:"name"`
		UserID string `json:"user_id"`
	}

	var data requestData

	// Parse request body into requestData struct
	if err := c.BodyParser(&data); err != nil {
		log.Println("Error parsing body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request!"})
	}

	if data.Name == "" || data.UserID == "" {
		log.Println("Error: missing fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	var existingOrganization models.Organization
	checkQuery := `SELECT organization_id FROM organizations WHERE LOWER(name) = LOWER($1)`
	err := db.QueryRow(context.Background(), checkQuery, data.Name).Scan(&existingOrganization.OrganizationID)
	if err == nil {
		// Folder with the same name exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Organization with the same name exists"})
	} else if err != pgx.ErrNoRows {
		log.Println("Error checking for existing folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing organization", "message": err.Error()})
	}

	// Initialize the Organization struct
	organization := models.Organization{
		OrganizationID: uuid.New().String(),
		Name:           data.Name,
		CreatedAt:      time.Now(),
	}

	query := "INSERT INTO organizations(organization_id, name, created_at) VALUES ($1, $2, $3);"

	_, err = db.Exec(
		context.Background(),
		query,
		organization.OrganizationID, organization.Name, organization.CreatedAt,
	)

	if err != nil {
		log.Println("Error creating organization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating organization", "message": err.Error()})
	}

	// Create the UserOrganization struct
	userOrganization := models.UserOrganization{
		ID:             uuid.New().String(),
		UserID:         data.UserID,
		OrganizationID: organization.OrganizationID,
		Role:           "creator",
		CreatedAt:      time.Now(),
	}

	query = `
		INSERT INTO userorganizations(id, user_id, organization_id, role, created_at)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err = db.Exec(
		context.Background(),
		query,
		userOrganization.ID, userOrganization.UserID, userOrganization.OrganizationID, userOrganization.Role, userOrganization.CreatedAt,
	)

	if err != nil {
		log.Println("Error creating user organization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user organization", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Organization created successfully!"})
}

func UpdateOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("id")

	if organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organizationId missing"})
	}

	var organization models.Organization

	if err := c.BodyParser(&organization); err != nil {
		log.Println("Error parsing body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request!"})
	}

	if organization.OrganizationID == "" || organization.Name == "" {
		log.Println("Error, missing fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	query := `
		UPDATE organizations
		SET name = $1
		WHERE organization_id = $2;
	`

	commandTag, err := db.Exec(
		context.Background(),
		query,
		organization.Name, organizationId,
	)

	if err != nil {
		log.Println("Error updating organization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating organization", "message": err.Error()})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Organization not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Organization updated successfully!"})
}

/**
	* DeleteOrganization - deletes an organization
	* @c - fiber
	* @db - database
	* Return - error or Ok
*/
func DeleteOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	organization_id := c.Params("organization_id")

	if organization_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Organization ID missing"})
	}

	query := "DELETE FROM organizations WHERE organization_id = $1;"

	_, err := db.Exec(
		context.Background(),
		query,
		organization_id,
	)

	if err != nil {
		log.Println("Error deleting organization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting organization", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "Organization deleted!"})
}

// Function to get all organizations for a specific user
func GetOrganizations(c *fiber.Ctx, db *pgxpool.Pool) error {
	userId := c.Params("user_id")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing user_id in params"})
	}

	query := `
		SELECT org.*
		FROM organizations org
		JOIN userorganizations uo ON org.organization_id = uo.organization_id 
		WHERE uo.user_id = $1;
	`

	rows, err := db.Query(
		context.Background(),
		query,
		userId,
	)

	if err != nil {
		log.Println("Error fetching organizations: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching organizations", "message": err.Error()})
	}
	defer rows.Close()

	var organizations []models.Organization
	for rows.Next() {
		var organization models.Organization
		if err := rows.Scan(
			&organization.OrganizationID, 
			&organization.Name, 
			&organization.CreatedAt,
			); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		organizations = append(organizations, organization)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(organizations)
}

//Function to fetch an organization
func GetOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	organization_id := c.Params("organization_id")
	var organization models.Organization

	if organization_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing organization_id"})
	}

	query := "SELECT organization_id, name, created_at FROM organizations WHERE organization_id = $1;"

	err := db.QueryRow(
		context.Background(),
		query,
		organization_id,
	).Scan(&organization.OrganizationID, &organization.Name, &organization.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Error fetching the organization", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(organization)
}

func RegisterOrganizationRoutes(app *fiber.App, db *pgxpool.Pool) {
	organizationGroup := app.Group("/organization")

	organizationGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateOrganization(c, db)
	})
	organizationGroup.Put("/update/:id", func(c *fiber.Ctx) error {
		return UpdateOrganization(c, db)
	})
	organizationGroup.Get("/fetch/specific/:organization_id", func(c *fiber.Ctx) error {
		return GetOrganization(c, db)
	})
	organizationGroup.Get("/fetch/all/:user_id", func(c *fiber.Ctx) error {
		return GetOrganizations(c, db)
	})
	organizationGroup.Delete("/delete/:organization_id", func(c *fiber.Ctx) error {
		return DeleteOrganization(c, db)
	})
}
