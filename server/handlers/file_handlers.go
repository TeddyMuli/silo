package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"server/spaces"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//Function to get files
func GetFiles(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("organization_id")
	folderId := c.Params("folder_id")

	if organizationId == ""  {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	var query string
	var rows pgx.Rows
	var err error

	if organizationId != "" && folderId == "" {
		query = `
		SELECT *
		FROM files 
		WHERE organization_id = $1 AND folder_id IS NULL AND deleted = false;
	`

		rows, err = db.Query(
			context.Background(),
			query,
			organizationId,
		)
	} else if organizationId != "" && folderId != "" {
		query =`
			SELECT *
			FROM files 
			WHERE folder_id = $1 AND deleted = false;
		`		

		rows, err = db.Query(
			context.Background(),
			query,
			folderId,
		)
	}

	if err != nil {
		log.Println("Error fetching files: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching files", "message": err.Error()})
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Name, &file.FolderID, &file.FilePath, &file.FileSize, &file.CreatedAt, &file.UpdatedAt, &file.OrganizationID, &file.Deleted, &file.DeletedAt); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(files)
}

//Function to get files
func GetFile(c *fiber.Ctx, db *pgxpool.Pool) error {
	fileId := c.Params("file_id")
	var file models.File

	if fileId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File id missing"})
	}

	query := `
		SELECT *
		FROM files 
		WHERE id = $1 AND deleted = false;
	`

	err := db.QueryRow(
		context.Background(),
		query,
		fileId,
	).Scan(
		&file.ID,
		&file.Name,
		&file.FolderID,
		&file.FilePath,
		&file.FileSize,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.OrganizationID,
		&file.Deleted,
		&file.DeletedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : "Error fetching file", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(file)
}

// CreateFile creates a new file
func CreateFile(c *fiber.Ctx, db *pgxpool.Pool) error {
	var file models.File
	if err := c.BodyParser(&file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	file.ID = uuid.New().String()
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	file.Deleted = false

	if file.Name == ""  {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Missing file name"})
	}

	if file.FolderID != nil && *file.FolderID == "null" {
		file.FolderID = nil
	}

	// Check if a file with the same name already exists
	/** 
	var existingFile models.File
	checkQuery := `
		SELECT id FROM files
		WHERE LOWER(name) = LOWER($1)
		AND organization_id = $2
		AND (folder_id = $3 OR (folder_id IS NULL AND $3 IS NULL))
		AND deleted = false
	`
	err := db.QueryRow(
		context.Background(),
		checkQuery,
		file.Name,
		file.OrganizationID,
		file.FolderID,
	).Scan(&existingFile.ID)

	if err == nil {
		// File with the same name exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "File with the same name already exists"})
	} else if err != pgx.ErrNoRows {
		log.Println("Error checking for existing file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing file", "message": err.Error()})
	}
	*/

	query := `
		INSERT INTO files
		(id, name, folder_id, file_path, file_size, created_at, updated_at, organization_id, deleted)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	_, err := db.Exec(
		context.Background(),
		query,
		file.ID, file.Name, file.FolderID, file.FilePath, file.FileSize, file.CreatedAt, file.UpdatedAt, file.OrganizationID, file.Deleted,
	)

	if err != nil {
		log.Println("Error creating file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating file", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(file)
}

// UpdateFile updates an existing file
func UpdateFile(c *fiber.Ctx, db *pgxpool.Pool) error {
	fileId := c.Params("id")
	var file models.File

	if fileId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fileId missing"})
	}

	if err := c.BodyParser(&file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	file.UpdatedAt = time.Now()

	// Dynamic query construction based on non-empty values
	var updateFields []string
	var args []interface{}
	argIndex := 1

	if file.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, file.Name)
		argIndex++
	}
	if *file.FolderID != "" {
		updateFields = append(updateFields, fmt.Sprintf("folder_id = $%d", argIndex))
		args = append(args, file.FolderID)
		argIndex++
	}
	if file.FilePath != "" {
		updateFields = append(updateFields, fmt.Sprintf("file_path = $%d", argIndex))
		args = append(args, file.FilePath)
		argIndex++
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, file.UpdatedAt)
	argIndex++

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	args = append(args, fileId)

	query := fmt.Sprintf(`
		UPDATE files
		SET %s
		WHERE id = $%d AND deleted = false;
	`, strings.Join(updateFields, ", "), argIndex)

	_, err := db.Exec(
		context.Background(),
		query,
		args...,
	)

	if err != nil {
		log.Println("Error updating file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating file", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(file)
}

// MoveFileTrash marks a file as deleted and sets a deletion timestamp
func MoveFileTrash(c *fiber.Ctx, db *pgxpool.Pool) error {
	fileId := c.Params("file_id")

	if fileId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fileId missing"})
	}

	deletedAt := time.Now().Add(30 * 24 * time.Hour) // 30 days from now

	query := "UPDATE files SET deleted = true, deleted_at = $1 WHERE id = $2 AND deleted = false;"

	_, err := db.Exec(
		context.Background(), 
		query, 
		deletedAt, fileId,
	)

	if err != nil {
		log.Println("Error deleting file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting file", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File marked for deletion"})
}

func DeleteFile(c *fiber.Ctx, db *pgxpool.Pool) error {
	fileId := c.Params("file_id")
	var file models.File

	if fileId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fileId missing"})
	}

	err := db.QueryRow(
		context.Background(), 
		"SELECT file_path FROM files WHERE id = $1",
		fileId,
	).Scan(&file.FilePath)

	if err != nil {
		log.Println("Error retrieving file_path", err)
		return fmt.Errorf("failed to retrieve file_path: %w", err)
	}

	err = spaces.DeleteFile(file.FilePath)
	if err != nil {
    log.Println("Error deleting file from Spaces:", err)
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting file from Spaces"})
	}

	query := "DELETE FROM files WHERE id = $1 AND deleted = true;"

	_, err = db.Exec(
		context.Background(), 
		query, 
		fileId,
	)

	if err != nil {
		log.Println("Error deleting file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting file", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully!"})
}

// RestoreFile updates a the deleted field to false and rstores the file to view
func RestoreFile(c *fiber.Ctx, db *pgxpool.Pool) error {
	fileId := c.Params("file_id")

	if fileId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fileId missing"})
	}

	query := "UPDATE files SET deleted = false WHERE id = $1 AND deleted = true;"

	_, err := db.Exec(
		context.Background(),
		query,
		fileId,
	)

	if err != nil {
		log.Println("Error restoring file: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error restoring file", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File restored"})
}

// DeleteExpiredFiles deletes files where the current date is past the deletedAt date
func DeleteExpiredFiles(db *pgxpool.Pool) {
	query := "DELETE FROM files WHERE deleted = true AND deleted_at < $1;"
	_, err := db.Exec(
		context.Background(),
		query,
		time.Now(),
	)

	if err != nil {
		log.Println("Error deleting expired files: ", err)
	} else {
		log.Println("Expired files deleted successfully")
	}
}

// GetDeletedFiless retrieves all files where deleted = true
//for a particular user
func GetDeletedFiles(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("organization_id")

	if organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organizationId missing"})
	}

	query := `
		SELECT id, name, folder_id, file_path, file_size, created_at, updated_at, organization_id, deleted_at
		FROM files
		WHERE deleted = true AND organization_id = $1 AND folder_id IS NULL;
	`

	rows, err := db.Query(context.Background(), query, organizationId)

	if err != nil {
		log.Println("Error querying deleted files: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying deleted files", "message": err.Error()})
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(&file.ID, &file.Name, &file.FolderID, &file.FilePath, &file.FileSize, &file.CreatedAt, &file.UpdatedAt, &file.OrganizationID, &file.DeletedAt)
		if err != nil {
			log.Println("Error scanning file: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning file", "message": err.Error()})
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating over rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(files)
}

func RegisterFileRoutes(app *fiber.App, db *pgxpool.Pool) {
	// File routes
	fileGroup := app.Group("/file")

	fileGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateFile(c, db)
	})
	fileGroup.Put("/update/:id", func(c *fiber.Ctx) error {
		return UpdateFile(c, db)
	})
	fileGroup.Get("/fetch/all/:organization_id/:folder_id?", func(c *fiber.Ctx) error {
		return GetFiles(c, db)
	})
	fileGroup.Get("/fetch/specific/:file_id", func(c *fiber.Ctx) error {
		return GetFile(c, db)
	})
	fileGroup.Get("/fetch/deleted/:organization_id", func(c *fiber.Ctx) error {
		return GetDeletedFiles(c, db)
	})
	fileGroup.Delete("/delete/permanent/:file_id", func(c *fiber.Ctx) error {
		return DeleteFile(c, db)
	})
	fileGroup.Put("/delete/:file_id", func(c *fiber.Ctx) error {
		return MoveFileTrash(c, db)
	})
	fileGroup.Put("/restore/:file_id", func(c *fiber.Ctx) error {
		return RestoreFile(c, db)
	})
}
