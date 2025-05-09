package db

import (
	"auth/config"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"
)

type Postgres struct {
	DB  *pgxpool.Pool
	log *zap.Logger
	ctx context.Context
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
	pgError    error
)

func NewPG(cfg *config.Config, log *zap.Logger, ctx context.Context) (*Postgres, error) {
	pgOnce.Do(func() {
		connString := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_min_conns=%d",
			cfg.Database.User, cfg.Database.Password,
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
			cfg.Database.SSLMode,
			cfg.Database.MaxConn, cfg.Database.MinConn,
		)

		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			pgError = fmt.Errorf("failed to create connection pool: %w", err)
			log.Error("Failed to create connection pool", zap.Error(err))
			return
		}

		if err = db.Ping(ctx); err != nil {
			pgError = fmt.Errorf("database ping failed: %w", err)
			log.Error("ManagementDatabase ping failed", zap.Error(err))
			db.Close()
			return
		}

		pgInstance = &Postgres{
			DB:  db,
			log: log,
			ctx: ctx,
		}

		migrateString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)

		m, err := migrate.New("file://./migrations", migrateString)
		if err != nil {
			pgError = fmt.Errorf("failed to init migrations: %w", err)
			log.Error("Failed to initialize migrations", zap.Error(err))
			return
		}

		log.Info("Applying database migrations")
		err = m.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			log.Info("No new migrations to apply")
		case err != nil && err.Error() == "Dirty database version 1. Fix and force version.":
			log.Warn("ManagementDatabase is dirty, forcing version 1")
			if forceErr := m.Force(1); forceErr != nil {
				pgError = fmt.Errorf("failed to force version: %w", forceErr)
				log.Error("Failed to force migration version", zap.Error(forceErr))
				return
			}
			if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				pgError = fmt.Errorf("migration failed after forcing version: %w", err)
				log.Error("Migration failed after forcing version", zap.Error(err))
				return
			}
		case err != nil:
			pgError = fmt.Errorf("migration failed: %w", err)
			log.Error("Migration failed", zap.Error(err))
			return
		}

		log.Info("ManagementDatabase connection established and migrations applied")
	})

	return pgInstance, pgError
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}
