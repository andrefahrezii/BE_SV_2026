package auth

import (
	"database/sql"
	"os"

	"github.com/you/sharing-vision-backend-v2/internal/config"
)

func RunMigrations(dbConn *sql.DB) error {
	if dbConn == nil {
		return nil
	}
	files := []string{
		"migrations/001_initial_schema.sql",
		"migrations/002_categories.sql",
		"migrations/003_fix_audit_log_action.sql",
	}
	for _, f := range files {
		schema, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := dbConn.Exec(string(schema)); err != nil {
			return err
		}
	}
	config.Log.Info("migrations applied")
	return nil
}
