package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thejixer/jixifood/shared/models"
)

func (s *PostgresStore) SeedDB() error {
	fmt.Println("seeeding the database initiating")
	var seeded bool
	err := s.db.QueryRow("SELECT EXISTS (SELECT 1 FROM seeded WHERE name = 'auth_service')").Scan(&seeded)
	if err != nil {
		return fmt.Errorf("failed to check seeding status: %w", err)
	}

	if seeded {
		log.Println("Auth service already seeded.")
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO roles (name, description)
		VALUES ('manager', 'Manager role with full access'),
					 ('operator', 'operator role with limited access'),
				   ('delivery', 'Dellivery agent'),
					 ('customer', 'customer')
		ON CONFLICT (name) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	roleIDs := make(map[string]int64)
	rows, err := tx.Query(`SELECT id, name FROM roles WHERE name IN ('manager', 'operator', 'delivery', 'customer')`)
	if err != nil {
		return fmt.Errorf("failed to retrieve role IDs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.Id, &role.Name); err != nil {
			return fmt.Errorf("failed to scan role: %w", err)
		}
		roleIDs[role.Name] = role.Id
	}

	// make the main manager
	phoneNumber := os.Getenv("MAIN_MANAGER_PHONENUMBER")
	lastInsertId := 0
	insertErr := tx.QueryRow(
		`	INSERT INTO USERS (phone_number, status, role_id, createdAt)
			VALUES ($1, $2, $3, $4) RETURNING id`,
		phoneNumber,
		"complete",
		roleIDs["manager"],
		time.Now().UTC(),
	).Scan(&lastInsertId)

	if insertErr != nil || lastInsertId == 0 {
		return fmt.Errorf("could not insert the main manager into the database")
	}

	// Define permissions and their role assignments
	permissions := map[string][]string{
		"manage_user":          {"manager"},
		"manage_menu":          {"manager", "operator"},
		"view orders":          {"manager", "operator"},
		"manage_orders":        {"manager", "operator"},
		"assign_orders":        {"manager", "operator"},
		"mark_as_delivered":    {"delivery"},
		"view_delivery_status": {"manager", "operator"},
	}

	// Insert permissions and assign them to roles
	permissionIDs := make(map[string]int)
	for permission, roles := range permissions {
		// Insert permission into the database
		_, err := tx.Exec(`INSERT INTO permissions (name, description) VALUES ($1, '') ON CONFLICT (name) DO NOTHING`, permission)
		if err != nil {
			return fmt.Errorf("failed to seed permission %s: %w", permission, err)
		}

		// Retrieve the permission ID
		var permissionID int
		err = tx.QueryRow(`SELECT id FROM permissions WHERE name = $1`, permission).Scan(&permissionID)
		if err != nil {
			return fmt.Errorf("failed to retrieve ID for permission %s: %w", permission, err)
		}
		permissionIDs[permission] = permissionID

		// Assign permission to the specified roles
		for _, role := range roles {
			roleID, exists := roleIDs[role]
			if !exists {
				return fmt.Errorf("role %s not found", role)
			}

			_, err := tx.Exec(`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT (role_id, permission_id) DO NOTHING`, roleID, permissionID)
			if err != nil {
				return fmt.Errorf("failed to assign permission %s to role %s: %w", permission, role, err)
			}
		}
	}

	// Mark the auth service as seeded
	_, err = tx.Exec("INSERT INTO seeded (name) VALUES ('auth_service')")
	if err != nil {
		return fmt.Errorf("failed to update seeded table: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Auth service seeding completed.")
	return nil

}
