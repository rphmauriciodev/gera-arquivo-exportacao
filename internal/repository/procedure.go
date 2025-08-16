package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetProcedureCode(db *sql.DB, dbType, procName string) (string, error) {
	var query string

	switch dbType {
	case "POSTGRES":
		query = `
		SELECT pg_get_functiondef(p.oid)
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		WHERE p.proname = $1
		LIMIT 1;
	`
	case "SQLSERVER":
		query = `
		SELECT sm.definition
		FROM sys.sql_modules sm
		INNER JOIN sys.objects so ON sm.object_id = so.object_id
		WHERE so.name = @p1
		`
	}

	var code sql.NullString
	err := db.QueryRow(query, procName).Scan(&code)
	if err != nil {
		return "", err
	}
	if !code.Valid {
		fmt.Printf("procedure %s não encontrada", procName)
		return "", nil
	}
	return strings.TrimSpace(code.String), nil
}
