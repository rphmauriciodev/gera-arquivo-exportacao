package config

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func ConnectDb(cfg *Config) (*sql.DB, error) {
	if cfg.Server.Server == "" || cfg.Server.User == "" || cfg.Server.Password == "" || cfg.Server.Database == "" {
		return nil, errors.New("configuração inválida: servidor, usuário, senha e banco são obrigatórios")
	}

	connString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;database=%s",
		cfg.Server.Server, cfg.Server.User, cfg.Server.Password, cfg.Server.Database,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco: %w", err)
	}

	return db, nil
}
