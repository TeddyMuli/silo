package handlers

import (
	"context"
	"log"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

/**
* CreateUserOrganization - Creates a user organization
* @c: fiber context
* db: database
* Return: Error or Ok
 */
func CreateUserOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	var userOrganization models.UserOrganization

	if err := c.BodyParser(&userOrganization); err != nil {
		log.Println("Error parsing body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request!"})
	}

	if userOrganization.OrganizationID == "" || userOrganization.UserID == "" || userOrganization.Role == "" {
		log.Println("Error, missing fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	userOrganization.ID = uuid.New().String()

	query := "INSERT INTO userorganizations (id, user_id, organization_id, role) VALUES ($1, $2, $3, $4);"

	_, err := db.Exec(
		context.Background(),
		query,
		userOrganization.ID, userOrganization.UserID, userOrganization.OrganizationID, userOrganization.Role,
	)

	if err != nil {
		log.Println("Error creating user organization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating organization", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User organization created successfully!"})
}

/**
	* CreateUserOrganization - Creates a user organization
	* @c: fiber context
	* db: database
	* Return: Error or Ok
*/
func UpdateUserOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	userId := c.Params("user_id")
	organizationId := c.Params("organization_id")
	
	if userId == "" || organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	var userOrganization models.UserOrganization

	if err := c.BodyParser(&userOrganization); err != nil {
		log.Println("Error parsing body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request!"})
	}

	if userOrganization.ID == "" || userOrganization.Role == "" {
		log.Println("Error, missing fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	query := `
		UPDATE userorganizations
		SET role = $1
		WHERE user_id = $2 AND organization_id = $3;
	`

	commandTag, err := db.Exec(
		context.Background(),
		query,
		userOrganization.Role, userId, organizationId,
	)

	if err != nil {
		log.Println("Error updating User Oorganization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating userorganization", "message": err.Error()})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User Organization not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User Organization updated successfully!"})
}

func DeleteUserOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	user_id := c.Params("user_id")
	organization_id := c.Params("organization_id")

	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "USer ID missing"})
	}

	query := "DELETE FROM userorganizations WHERE user_id = $1 AND organization_id = $2;"

	_, err := db.Exec(
		context.Background(),
		query,
		user_id, organization_id,
	)

	if err != nil {
		log.Println("Error deleting user organization")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting user organization", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User Organization deleted successfully!"})
}

//Function to fetch organizations a user is part of
func GetUserOrganizations(c *fiber.Ctx, db *pgxpool.Pool) error {
	user_id := c.Params("user_id")

	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id missing in the request"})
	}

	query := `
		SELECT 
			o.organization_id,
			o.name,
			o.created_at,
			uo.role,
			uo.created_at AS user_organization_created_at
		FROM 
			UserOrganizations uo
		JOIN 
			Organizations o ON uo.organization_id = o.organization_id
		WHERE 
			uo.user_id = $1;
	`

	rows, err := db.Query(
		context.Background(),
		query,
		user_id,
	)

	if err != nil {
		log.Println("Error fetching organization", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching organization", "message": err.Error()})
	}
	defer rows.Close()

	var organizations []map[string]interface{}
	var organization models.Organization
	var userOrganization models.UserOrganization

	for rows.Next() {
		if err := rows.Scan(&organization.OrganizationID, &organization.Name, &organization.CreatedAt, &userOrganization.Role, &userOrganization.CreatedAt); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}

		organization := map[string]interface{}{
				"organization_id":               organization.OrganizationID,
				"name":                          organization.Name,
				"created_at":                    organization.CreatedAt,
				"role":                          userOrganization.Role,
				"user_organization_created_at":  userOrganization.CreatedAt,
		}
		organizations = append(organizations, organization)
}
	return c.Status(fiber.StatusOK).JSON(organizations)
}

func RegisterUserOrganizationRoutes(app *fiber.App, db *pgxpool.Pool) {
	userOrganizatonGroup := app.Group("/user_organization")

	userOrganizatonGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateUserOrganization(c, db)
	})
	userOrganizatonGroup.Put("/update/:user_id/:organization_id", func(c *fiber.Ctx) error {
		return UpdateUserOrganization(c, db)
	})
	userOrganizatonGroup.Get("/fetch/:user_id", func(c *fiber.Ctx) error {
		return GetUserOrganizations(c, db)
	})
	userOrganizatonGroup.Delete("/delete/:user_id/:organization_id", func(c *fiber.Ctx) error {
		return DeleteUserOrganization(c, db)
	})
}
