package handlers

import (
	"context"
	"log"

	"server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePlan(c *fiber.Ctx, db *pgxpool.Pool) error {
	var plan models.Plan

	if err := c.BodyParser(&plan); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	if plan.Name == "" || plan.Description == "" || plan.Interval == "" || plan.Price == 0 || plan.StorageLimit == 0 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing values"})
	}

	plan.ID = uuid.New().String()

	query := "INSERT INTO plans (id, name, description, price, interval, storage_limit) VALUES ($1, $2, $3, $4, $5, $6);"

	_, err := db.Exec(
		context.Background(),
		query,
		plan.ID, plan.Name, plan.Description, plan.Price, plan.Interval, plan.StorageLimit,
	)

	if err != nil {
		log.Println("Error creating plan:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : "Error creating plan", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Plan created successfully!"})
}

func UpdatePlan(c *fiber.Ctx, db *pgxpool.Pool) error {
	planId := c.Params("id")

	if planId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "planId missing"})
	}

	var plan models.Plan

	if err := c.BodyParser(&plan); err != nil {
		log.Println("Error parsing", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error" : "Error parsing request, Invalid data!"})
	}

	query := `
		UPDATE plans
		SET name=$1, description=$2, price=$3, interval=$4, storage_limit=$5
		WHERE id=$6;
	`

	_, err := db.Exec(
		context.Background(),
		query,
		plan.Name, plan.Description, plan.Price, plan.Interval, plan.StorageLimit, planId,
	)

	if err != nil {
		log.Println("Error updating plan: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating plan", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Plan updated!"})
}

func DeletePlan(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := "DELETE FROM plans WHERE id=$1;"

	_, err := db.Exec(context.Background(), query, id)

	if err != nil {
		log.Println("Error deleting plan:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting plan", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Plan deleted successfully!"})

}

func GetPlan(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	var plan models.Plan

	query := "SELECT id, name, description, price, interval, storage_limit FROM plans WHERE id=$1;"

	err := db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.Price, &plan.Interval, &plan.StorageLimit,
	)

	if err != nil {
		log.Println("Error retrieving plan:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving plan", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(plan)
}

func GetPlans(c *fiber.Ctx, db *pgxpool.Pool) error {
	var plans []models.Plan

	query := "SELECT * FROM plans;"

	rows, err := db.Query(
		context.Background(),
		query,
	)

	if err != nil {
		log.Println("Error retrieving plans:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving plans", "message": err.Error()})
	}

	for rows.Next() {
		var plan models.Plan
		if err := rows.Scan(&plan.ID, &plan.Name, &plan.Description, &plan.Price, &plan.Interval, &plan.StorageLimit); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		plans = append(plans, plan)
	}

	return c.Status(fiber.StatusOK).JSON(plans)
}

func RegisterPlanRoutes(app *fiber.App, db *pgxpool.Pool) {
	planGroup := app.Group("/plan")

	planGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreatePlan(c, db)
	})
	planGroup.Put("/update/:id", func(c *fiber.Ctx) error {
		return UpdatePlan(c, db)
	})
	planGroup.Get("/fetch/specific/:id", func(c *fiber.Ctx) error {
		return GetPlan(c, db)
	})
	planGroup.Get("/fetch/all", func(c *fiber.Ctx) error {
		return GetPlans(c, db)
	})
	planGroup.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return DeletePlan(c, db)
	})
}
