package postgresql

import (
	"catalog/internal/config"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const UnableToConnectDatabase = "Unable to connect to database: "

type Storage struct {
	DB  *sql.DB
	log *slog.Logger
}

func ConnectAndNew(log *slog.Logger, cfg *config.DatabaseConfig) (*Storage, error) {
	const op = "database.postgresql.ConnectAndNew"

	log = log.With(
		slog.String("op", op),
	)

	dsn := getDSN(cfg)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error(UnableToConnectDatabase, "error", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error(UnableToConnectDatabase, "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(80)

	storage := &Storage{
		DB:  db,
		log: log,
	}

	return storage, nil
}

func NewRep(db *sql.DB, log *slog.Logger) *Storage {
	return &Storage{
		DB:  db,
		log: log,
	}
}

func getDSN(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s target_session_attrs=read-write", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
}

func (s *Storage) Stop() {
	s.DB.Close()
}
