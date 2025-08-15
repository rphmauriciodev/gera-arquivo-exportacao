package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetProcedureCode(db *sql.DB, procName string) (string, error) {
	query := `
		SELECT sm.definition
		FROM sys.sql_modules sm
		INNER JOIN sys.objects so ON sm.object_id = so.object_id
		WHERE so.name = @p1
	`
	var code sql.NullString
	err := db.QueryRow(query, procName).Scan(&code)
	if err != nil {
		return "", err
	}
	if !code.Valid {
		return "", fmt.Errorf("procedure %s não encontrada", procName)
	}
	return strings.TrimSpace(code.String), nil
}
