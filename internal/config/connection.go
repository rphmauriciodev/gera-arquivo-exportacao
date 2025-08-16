package config

import (
	"context"
	"database/sql"
	"errors"
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

func (p *PostgresDB) Connect(cfg *Config) (*sql.DB, error) {
	if cfg.Server.Server == "" || cfg.Server.User == "" || cfg.Server.Password == "" || cfg.Server.Database == "" {
		return nil, errors.New("configuração inválida: host, usuário, senha e banco são obrigatórios")
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.Server.User, cfg.Server.Password, cfg.Server.Server, cfg.Server.Database,
	)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com PostgreSQL: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("falha ao conectar no PostgreSQL: %w", err)
	}

	return db, nil
}
