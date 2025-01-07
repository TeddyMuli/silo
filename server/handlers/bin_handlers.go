package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func DeleteExpiredItems(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("organization_id")

	if organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing organizationId"})
	}
	// Delete from folders
	folderQuery := "DELETE FROM folders WHERE organization_id = $1 AND deleted = true;"
	_, err := db.Exec(
		context.Background(),
		folderQuery,
		organizationId,
	)

	if err != nil {
		log.Println("Error deleting folders: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error clearing bin"})
	}

	// Delete from files
	fileQuery := "DELETE FROM files WHERE organization_id = $1 AND deleted = true;"
	_, err = db.Exec(
		context.Background(),
		fileQuery,
		organizationId,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error clearing bin"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Recycle bin cleared successfully!"})
}

func RegisterBinRoutes (app *fiber.App, db *pgxpool.Pool) {
	app.Delete("/bin/empty/:organization_id", func(c *fiber.Ctx) error {
		return DeleteExpiredItems(c, db)
	})
}
