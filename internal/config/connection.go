package config

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type IDatabase interface {
	Connect(cfg *Config) (*sql.DB, error)
}

type SQLServerDB struct{}
type PostgresDB struct{}

func (s *SQLServerDB) Connect(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", cfg.Server.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco: %w", err)
	}

	return db, nil
}

func (p *PostgresDB) Connect(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Server.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com PostgreSQL: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("falha ao conectar no PostgreSQL: %w", err)
	}

	return db, nil
}
