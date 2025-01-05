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

// Function to get folders
func GetFolders(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("organization_id")

	if organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organizationId missing"})
	}

	query := `
		SELECT id, name, organization_id, COALESCE(parent_folder_id, NULL), created_at, updated_at, deleted, deleted_at
		FROM folders 
		WHERE organization_id = $1 AND deleted = false;
	`

	rows, err := db.Query(
		context.Background(),
		query,
		organizationId,
	)

	if err != nil {
		log.Println("Error fetching folders: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching folders", "message": err.Error()})
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		if err := rows.Scan(
			&folder.ID, 
			&folder.Name, 
			&folder.OrganizationID, 
			&folder.ParentFolderID, 
			&folder.CreatedAt, 
			&folder.UpdatedAt, 
			&folder.Deleted, 
			&folder.DeletedAt); 
			err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		if folder.ParentFolderID == nil {
			emptyStr := ""
			folder.ParentFolderID = &emptyStr
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(folders)
}

// Function to get child folders
func GetChildFolders(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("organization_id")
	parentFolderID := c.Params("parent_folder_id")

	if parentFolderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "parent_folder_id missing"})
	}

	var query string
	var rows pgx.Rows
	var err error

	if parentFolderID == "root" {
		// If parentFolderID is "root", fetch folders where parent_folder_id is NULL
		query := `
			SELECT id, name, organization_id, parent_folder_id, created_at, updated_at, deleted, deleted_at
			FROM folders 
			WHERE organization_id = $1 AND parent_folder_id IS NULL AND deleted = false;
		`
		rows, err = db.Query(context.Background(), query, organizationId)
	} else {
		// Otherwise, fetch folders where parent_folder_id matches the provided value
		query = `
			SELECT id, name, organization_id, parent_folder_id, created_at, updated_at, deleted, deleted_at
			FROM folders 
			WHERE organization_id = $1 AND parent_folder_id = $2 AND deleted = false;
		`
		rows, err = db.Query(context.Background(), query, organizationId, parentFolderID)
	}

	if err != nil {
		log.Println("Error fetching folders: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching folders", "message": err.Error()})
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		if err := rows.Scan(
			&folder.ID, 
			&folder.Name, 
			&folder.OrganizationID, 
			&folder.ParentFolderID, 
			&folder.CreatedAt, 
			&folder.UpdatedAt, 
			&folder.Deleted, 
			&folder.DeletedAt); 
			err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row", "message": err.Error()})
		}
		if folder.ParentFolderID == nil {
			emptyStr := ""
			folder.ParentFolderID = &emptyStr
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(folders)
}


// Function to fetch a specific folder
func GetFolder(c *fiber.Ctx, db *pgxpool.Pool) error {
	var folder models.Folder
	folderId := c.Params("folder_id")

	if folderId == ""  {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "folderId missing"})
	}

	query := `
		SELECT
		id, name, organization_id, COALESCE(parent_folder_id, NULL), created_at, updated_at, deleted, deleted_at
		FROM folders
		WHERE id = $1;
	`

	err := db.QueryRow(
		context.Background(),
		query,
		folderId,
	).Scan(
		&folder.ID,
		&folder.Name,
		&folder.OrganizationID,
		&folder.ParentFolderID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.Deleted,
		&folder.DeletedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching folder", "message": err.Error()})
	}

	if folder.ParentFolderID == nil {
		emptyStr := ""
		folder.ParentFolderID = &emptyStr
	}

	return c.Status(fiber.StatusOK).JSON(folder)
}

// CreateFolder creates a new folder
func CreateFolder(c *fiber.Ctx, db *pgxpool.Pool) error {
	var folder models.Folder
	if err := c.BodyParser(&folder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	folder.ID = uuid.New().String()
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()
	folder.Deleted = false

	if folder.Name == "" || folder.OrganizationID == ""  {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Missing fields"})
	}

	// Check if a folder with the same name already exists
	/** 
	var existingFolder models.Folder
	checkQuery := `
		SELECT id FROM folders
		WHERE LOWER(name) = LOWER($1)
		AND organization_id = $2
		AND (parent_folder_id = $3 OR (parent_folder_id IS NULL AND $3 IS NULL)) 
		AND deleted = false
	`

	err := db.QueryRow(
		context.Background(),
		checkQuery,
		folder.Name, folder.OrganizationID, folder.ParentFolderID,
	).Scan(&existingFolder.ID)

	if err == nil {
		// Folder with the same name exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Folder with the same name exists"})
	} else if err != pgx.ErrNoRows {
		log.Println("Error checking for existing folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing folder", "message": err.Error()})
	}
	*/

	query := `
		INSERT INTO folders
		(id, name, organization_id, parent_folder_id, created_at, updated_at, deleted)
		VALUES
		($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := db.Exec(
		context.Background(),
		query,
		folder.ID, folder.Name, folder.OrganizationID, folder.ParentFolderID, folder.CreatedAt, folder.UpdatedAt, folder.Deleted,
	)

	if err != nil {
		log.Println("Error creating folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating folder", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Folder created!",
		"id": folder.ID,
	})
}

// UpdateFolder updates an existing folder
func UpdateFolder(c *fiber.Ctx, db *pgxpool.Pool) error {
	folderId := c.Params("folder_id")

	if folderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "folderId missing"})
	}

	var folder models.Folder
	if err := c.BodyParser(&folder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	folder.UpdatedAt = time.Now()

	var updateFields []string
	var args []interface{}
	argIndex := 1

	if folder.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, folder.Name)
		argIndex++
	}

	if folder.ParentFolderID != nil && *folder.ParentFolderID != "" {
		updateFields = append(updateFields, fmt.Sprintf("parent_folder_id = $%d", argIndex))
		args = append(args, *folder.ParentFolderID)
		argIndex++
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, folder.UpdatedAt)
	argIndex++

	query := fmt.Sprintf(`
		UPDATE folders
		SET %s
		WHERE id = $%d AND deleted = false;
	`, strings.Join(updateFields, ", "), argIndex)

	args = append(args, folderId)

	_, err := db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Println("Error updating folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating folder", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(folder)
}

// MoveFolderTrash marks a folder as deleted and sets a deletion timestamp
func MoveFolderTrash(c *fiber.Ctx, db *pgxpool.Pool) error {
	folderId := c.Params("folder_id")

	if folderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Folder id missing"})
	}

	deletedAt := time.Now().Add(30 * 24 * time.Hour) // 30 days from now

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error starting transaction", "message": err.Error()})
	}
	defer tx.Rollback(context.Background())

	folderQuery := "UPDATE folders SET deleted = true, deleted_at = $1 WHERE (id = $2 OR parent_folder_id = $2) AND deleted = false;"
	filesQuery := "UPDATE files SET deleted = true, deleted_at = $1 WHERE folder_id = $2 AND deleted = false;"

	// Delete a folder
	_, err = tx.Exec(context.Background(), folderQuery, deletedAt, folderId)
	if err != nil {
		log.Println("Error deleting folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting folder", "message": err.Error()})
	}

	// Delete files
	_, err = tx.Exec(context.Background(), filesQuery, deletedAt, folderId)
	if err != nil {
		log.Println("Error deleting files: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting files", "message": err.Error()})
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Error committing transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error committing transaction", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Folder and files marked for deletion"})
}

//DeleteFolder dletes a folder, its subfolders and files
func DeleteFolder(c *fiber.Ctx, db *pgxpool.Pool) error {
	folderId := c.Params("folder_id")
	var rows pgx.Rows

	if folderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Folder id missing"})
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error starting transaction", "message": err.Error()})
	}
	defer tx.Rollback(context.Background())

	folderQuery := "DELETE FROM folders WHERE (id =$1 OR parent_folder_id = $1) AND deleted = true"
	filesQuery := "DELETE FROM files WHERE folder_id =$1 AND deleted = true"

	// Delete a folder and subfolders
	_, err = tx.Exec(context.Background(), folderQuery, folderId)
	if err != nil {
		log.Println("Error deleting folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting folder", "message": err.Error()})
	}

	// Select file_paths of the slected files and delete them from spaces
	rows, err = db.Query(
		context.Background(),
		"SELECT file_path FROM files WHERE folder_id = $1",
		folderId,
	)

	if err != nil {
		log.Println("Error fetching file_paths")
		return fmt.Errorf("error fetching file paths: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.FilePath); err != nil {
			log.Println("Error scanning row: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning row"})
		}

		err = spaces.DeleteFile(file.FilePath)
		if err != nil {
			log.Println("Error deleting file from Spaces:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting file from Spaces"})
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating rows", "message": err.Error()})
	}

	// Delete files
	_, err = tx.Exec(context.Background(), filesQuery, folderId)
	if err != nil {
		log.Println("Error deleting files: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting files", "message": err.Error()})
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Error committing transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error committing transaction", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Folder and files deleted successfully"})
}

// RestoreFolder updates a the deleted field to false and rstores the folder to view
func RestoreFolder(c *fiber.Ctx, db *pgxpool.Pool) error {
	folderId := c.Params("folder_id")

	if folderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Folder id missing"})
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error starting transaction", "message": err.Error()})
	}
	defer tx.Rollback(context.Background())

	//Restore folder
	query := "UPDATE folders SET deleted = false WHERE (id = $1 OR parent_folder_id = $1) AND deleted = true;"
	_, err = db.Exec(
		context.Background(),
		query,
		folderId,
	)
	if err != nil {
		log.Println("Error restoring folder: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error restoring folder", "message": err.Error()})
	}

	// Update the files within the folders
	filesQuery := "UPDATE files SET deleted = false WHERE folder_id = $1 AND deleted = true;"
	_, err = tx.Exec(context.Background(), filesQuery, folderId)
	if err != nil {
		log.Println("Error restoring files: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error restoring files", "message": err.Error()})
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Error committing transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error committing transaction", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Folder restored"})
}

// DeleteExpiredFolders deletes folders where the current date is past the deletedAt date
func DeleteExpiredFolders(db *pgxpool.Pool) {
	query := "DELETE FROM folders WHERE deleted = true AND deleted_at < $1;"
	_, err := db.Exec(
		context.Background(),
		query,
		time.Now(),
	)

	if err != nil {
		log.Println("Error deleting expired folders: ", err)
	} else {
		log.Println("Expired folders deleted successfully")
	}
}

// GetDeletedFolders retrieves all folders where deleted = true
func GetDeletedFolders(c *fiber.Ctx, db *pgxpool.Pool) error {
	organizationId := c.Params("organization_id")

	if organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organizationId missing"})
	}

	query := `
		SELECT id, name, parent_folder_id, created_at, updated_at, deleted, deleted_at
		FROM folders
		WHERE deleted = true 
		AND organization_id = $1 AND parent_folder_id IS NULL;
	`
	rows, err := db.Query(context.Background(), query, organizationId)
	if err != nil {
		log.Println("Error querying deleted folders: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying deleted folders", "message": err.Error()})
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(&folder.ID, &folder.Name, &folder.ParentFolderID, &folder.CreatedAt, &folder.UpdatedAt, &folder.Deleted, &folder.DeletedAt)
		if err != nil {
			log.Println("Error scanning folder: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning folder", "message": err.Error()})
		}
		folders = append(folders, folder)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating over rows", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(folders)
}

func RegisterFolderRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Folder routes
	folderGroup := app.Group("/folder")

	folderGroup.Post("/create", func(c *fiber.Ctx) error {
		return CreateFolder(c, db)
	})
	folderGroup.Put("/update/:folder_id", func(c *fiber.Ctx) error {
		return UpdateFolder(c, db)
	})
	folderGroup.Put("/restore/:folder_id", func(c *fiber.Ctx) error {
		return RestoreFolder(c, db)
	})
	folderGroup.Put("/delete/:folder_id", func(c *fiber.Ctx) error {
		return MoveFolderTrash(c, db)
	})
	folderGroup.Delete("/delete/permanent/:folder_id", func(c *fiber.Ctx) error {
		return DeleteFolder(c, db)
	})
	folderGroup.Get("/fetch/all/:organization_id", func(c *fiber.Ctx) error {
		return GetFolders(c, db)
	})
	folderGroup.Get("/fetch/children/:organization_id/:parent_folder_id", func(c *fiber.Ctx) error {
		return GetChildFolders(c, db)
	})
	folderGroup.Get("/fetch/specific/:folder_id", func(c *fiber.Ctx) error {
		return GetFolder(c, db)
	})
	folderGroup.Get("/fetch/deleted/:organization_id", func(c *fiber.Ctx) error {
		return GetDeletedFolders(c, db)
	})
}
