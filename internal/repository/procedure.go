package repository

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func ReadProceduresFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var procs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			procs = append(procs, line)
		}
	}
	return procs, scanner.Err()
}

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
