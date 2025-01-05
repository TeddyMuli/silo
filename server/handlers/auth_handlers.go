package handlers

import (
	"context"
	"errors"
	"log"
	"os"

	//"os"
	"server/models"
	"server/otp"
	"server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Function for logging in
func Login(c *fiber.Ctx, db *pgxpool.Pool) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	var storedUser models.User
	err := db.QueryRow(
		context.Background(),
		"SELECT user_id, email, password, first_name, last_name FROM users WHERE email=$1",
		user.Email,
	).Scan(&storedUser.UserID, &storedUser.Email, &storedUser.Password, &storedUser.FirstName, &storedUser.LastName)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User doesn't exist!", "details": err.Error()})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user", "details": err.Error()})
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
			"details": err.Error(),
		})
	}

	_otp, err := otp.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating OTP"})
	}

	if err := otp.StoreOTP(storedUser.Email, _otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error storing OTP"})
	}

	// Send OTP email
	if err := otp.SendOTPEmail(storedUser.Email, _otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error sending OTP email",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OTP sent. Please verify."})
}

func VerifyOTP(c *fiber.Ctx, db *pgxpool.Pool) error {
	type OTPRequest struct {
		Email string `json:"email"`
		OTP    string `json:"otp"`
	}

	var data OTPRequest

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	var storedUser models.User
	err := db.QueryRow(
		context.Background(),
		"SELECT user_id, email, password, first_name, last_name FROM users WHERE email=$1",
		data.Email,
	).Scan(&storedUser.UserID, &storedUser.Email, &storedUser.Password, &storedUser.FirstName, &storedUser.LastName)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user", "details": err.Error()})
	}

	storedOTP, err := otp.GetStoredOTP(data.Email)
	if err != nil || storedOTP != data.OTP {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}

	token, err := utils.GenerateToken(storedUser.UserID, storedUser.Email, storedUser.FirstName, storedUser.LastName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token", "message": err.Error()})
	}

	secureCookie := os.Getenv("SECURE_COOKIE") == "true"

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		MaxAge:		int(time.Hour.Seconds() * 24 * 3),
		HTTPOnly: true,
		Secure:   secureCookie,
		SameSite: "None",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful!", "token": token})
}

// Function for registering a new user
func Register(c *fiber.Ctx, db *pgxpool.Pool) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		log.Println("Error parsing body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	if user.Email == "" || user.Password == "" {
		log.Println("Validation failed: missing fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password", "message": err.Error()})
	}

	user.UserID = uuid.New().String()

	err = db.QueryRow(
		context.Background(),
		`
			INSERT INTO
			users (user_id, first_name, last_name, email, phone_number, password, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (email) DO NOTHING
			RETURNING user_id
		`,
		user.UserID, user.FirstName, user.LastName, user.Email, user.PhoneNumber, string(hashedPassword), time.Now(),
	).Scan(&user.UserID)
	if err != nil {
		log.Println("Error registering user: ", err)
		if err == pgx.ErrNoRows {
			// No row was inserted due to conflict
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error registering user",
			"details": err.Error(),
		})
	}

	// Create a new organization
	var organization models.Organization
	organization.OrganizationID = uuid.New().String()
	organization.Name = user.FirstName + "'s Organization"
	organization.CreatedAt = time.Now()

	query := "INSERT INTO organizations (organization_id, name, created_at) VALUES ($1, $2, $3);"
	_, err = db.Exec(
		context.Background(),
		query,
		organization.OrganizationID, organization.Name, organization.CreatedAt,
	)
	if err != nil {
		log.Println("Error creating organization: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error creating organization",
			"message": err.Error(),
		})
	}

	// Link the user to the organization with a specific role
	var userOrganization models.UserOrganization
	userOrganization.ID = uuid.New().String()
	userOrganization.UserID = user.UserID
	userOrganization.OrganizationID = organization.OrganizationID
	userOrganization.Role = "creator"

	query = `
		INSERT INTO
		userorganizations (id, user_id, organization_id, role)
		VALUES ($1, $2, $3, $4);
	`
	_, err = db.Exec(
		context.Background(),
		query,
		userOrganization.ID, userOrganization.UserID, userOrganization.OrganizationID, userOrganization.Role,
	)
	if err != nil {
		log.Println("Error creating user organization: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error creating user organization",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":         "User registered and organization created successfully",
		"user_id":         user.UserID,
		"organization_id": organization.OrganizationID,
	})
}

// Function to delete a user, expects email in the request query
func DeleteUser(c *fiber.Ctx, db *pgxpool.Pool) error {
	email := c.Params("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email missing"})
	}

	query := "DELETE FROM users WHERE email = $1;"

	_, err := db.Exec(
		context.Background(),
		query,
		email,
	)

	if err != nil {
		log.Println("Error deleting user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting user", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted!"})
}

// Function to list all the users
func GetUsers(c *fiber.Ctx, db *pgxpool.Pool) error {
	query := `
		SELECT
		COALESCE(email, '') AS email,
		COALESCE(first_name, '') AS first_name,
		COALESCE(last_name, '') AS last_name 
		FROM users;
    `

	rows, err := db.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error listing users", "message": err.Error()})
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Email, &user.FirstName, &user.LastName); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning user data", "message": err.Error()})
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error with user data", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetOrganizationsCreatedByUser(c *fiber.Ctx, db *pgxpool.Pool) error {
	userId := c.Params("user_id")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing user_id"})
	}

	query := `
		SELECT org.*
		FROM organizations org
		JOIN userorganizations uo ON org.organization_id = uo.organization_id 
		WHERE uo.user_id = $1 AND uo.role = 'creator';
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

// Function to get a user's data, expects email in the request query
func GetUser(c *fiber.Ctx, db *pgxpool.Pool) error {
	user_id := c.Params("user_id")
	var user models.User

	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User Id missing"})
	}

	query := "SELECT first_name, last_name, phone_number FROM users WHERE user_id = $1;"
	err := db.QueryRow(
		context.Background(),
		query,
		user_id,
	).Scan(&user.FirstName, &user.LastName, &user.PhoneNumber)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching the user", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Function to update a user
func UpdateUser(c *fiber.Ctx, db *pgxpool.Pool) error {
	email := c.Params("email")
	var user models.User

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email missing"})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	query := `
        UPDATE users
        SET first_name = $1, last_name = $2, phone_number = $3
        WHERE email = $4;
    `

	commandTag, err := db.Exec(
		context.Background(),
		query,
		user.FirstName, user.LastName, user.PhoneNumber, email,
	)

	if err != nil {
		log.Println("Error updating user: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating user", "message": err.Error()})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

func RegisterAuthRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Authentication Routes
	authGroup := app.Group("/auth")

	authGroup.Post("/register", func(c *fiber.Ctx) error {
		return Register(c, db)
	})
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return Login(c, db)
	})
	authGroup.Post("/verify", func(c *fiber.Ctx) error {
		return VerifyOTP(c, db)
	})
	authGroup.Get("/fetch/all", func(c *fiber.Ctx) error {
		return GetUsers(c, db)
	})
	authGroup.Get("/fetch/specific/:user_id", func(c *fiber.Ctx) error {
		return GetUser(c, db)
	})
	authGroup.Get("/fetch/organizations/:user_id", func(c *fiber.Ctx) error {
		return GetOrganizationsCreatedByUser(c, db)
	})
	authGroup.Put("/update/:email", func(c *fiber.Ctx) error {
		return UpdateUser(c, db)
	})
	authGroup.Delete("/delete/:email", func(c *fiber.Ctx) error {
		return DeleteUser(c, db)
	})
}
