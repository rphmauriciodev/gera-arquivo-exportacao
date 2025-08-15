package db

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;database=%s",
		cfg.Server, cfg.Server.User, cfg.Server.Password, cfg.Server.Database,
	)
	return sql.Open("sqlserver", connString)
}

