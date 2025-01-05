package handlers

import (
	"context"
	"log"
	"time"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CreateSubscription inserts a new subscription into the database
func CreateSubscription(c *fiber.Ctx, db *pgxpool.Pool) error {
	var sub models.Subscription

	if err := c.BodyParser(&sub); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	sub.ID = uuid.New().String()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()

	query := `
		INSERT INTO subscriptions (id, organization_id, plan_id, status, stripe_subscription_id, start_date, end_date, paypal_subscription_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	_, err := db.Exec(
		context.Background(),
		query,
		sub.ID, sub.OrganizationID, sub.PlanID, sub.Status, sub.StripeSubscriptionID, sub.StartDate, sub.EndDate, sub.PaypalSubscriptionID, sub.CreatedAt, sub.UpdatedAt,
	)

	if err != nil {
		log.Println("Error creating subscription:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating subscription", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Subscription created successfully!"})
}

// GetSubscriptionByID retrieves a subscription by its ID
func GetSubscriptionByID(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	var sub models.Subscription

	query := `
		SELECT id, organization_id, plan_id, status, stripe_subscription_id, start_date, end_date, paypal_subscription_id, created_at, updated_at
		FROM subscriptions
		WHERE id = $1;
	`
	err := db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(&sub.ID, &sub.OrganizationID, &sub.PlanID, &sub.Status, &sub.StripeSubscriptionID, &sub.StartDate, &sub.EndDate, &sub.PaypalSubscriptionID, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		log.Println("Error retrieving subscription:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving subscription", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(sub)
}

// GetSubscriptionByID retrieves a subscription by its organization_id
func GetSubscriptionByOrganization(c *fiber.Ctx, db *pgxpool.Pool) error {
	organization_id := c.Params("organization_id")

	if organization_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organization_id missing"})
	}

	var sub models.Subscription

	query := `
		SELECT id, organization_id, plan_id, status, stripe_subscription_id, start_date, end_date, paypal_subscription_id, created_at, updated_at
		FROM subscriptions
		WHERE organization_id = $1;
	`
	err := db.QueryRow(
		context.Background(),
		query,
		organization_id,
	).Scan(&sub.ID, &sub.OrganizationID, &sub.PlanID, &sub.Status, &sub.StripeSubscriptionID, &sub.StartDate, &sub.EndDate, &sub.PaypalSubscriptionID, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		log.Println("Error retrieving subscription:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving subscription", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(sub)
}

// UpdateSubscription updates an existing subscription
func UpdateSubscription(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	var sub models.Subscription

	if err := c.BodyParser(&sub); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	sub.UpdatedAt = time.Now()

	query := `
		UPDATE subscriptions
		SET organization_id = $2, plan_id = $3, status = $4, stripe_subscription_id = $5, start_date = $6, end_date = $7, paypal_subscription_id = $8, updated_at = $9
		WHERE id = $1;
	`
	_, err := db.Exec(context.Background(), query, id, sub.OrganizationID, sub.PlanID, sub.Status, sub.StripeSubscriptionID, sub.StartDate, sub.EndDate, sub.PaypalSubscriptionID, sub.UpdatedAt)
	if err != nil {
		log.Println("Error updating subscription:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating subscription", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Subscription updated successfully!"})
}

// DeleteSubscription deletes a subscription by its ID
func DeleteSubscription(c *fiber.Ctx, db *pgxpool.Pool) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id missing"})
	}

	query := `
		DELETE FROM subscriptions
		WHERE id = $1;
	`
	_, err := db.Exec(context.Background(), query, id)
	if err != nil {
		log.Println("Error deleting subscription:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting subscription", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Subscription deleted successfully!"})
}

func RegisterSubscriptionRoutes(app *fiber.App, db *pgxpool.Pool) {
	subscriptionGroup := app.Group("/subscription")

	subscriptionGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateSubscription(c, db)
	})
	subscriptionGroup.Put("/update/:id", func(c *fiber.Ctx) error {
		return UpdateSubscription(c, db)
	})
	subscriptionGroup.Get("/fetch/subscription/:id", func(c *fiber.Ctx) error {
		return GetSubscriptionByID(c, db)
	})
	subscriptionGroup.Get("/fetch/organization/:organization_id", func(c *fiber.Ctx) error {
		return GetSubscriptionByOrganization(c, db)
	})
	subscriptionGroup.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return DeleteSubscription(c, db)
	})
}

