package migration

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"nutgaard/dora-metrics/internal/utils"
	"os"
	"path/filepath"
	"time"
)

type MigratorCtx struct {
	pool *pgxpool.Pool
}

type status = string

const (
	Success      status = "Success"
	Failure      status = "Failure"
	Pending      status = "Pending"
	NotCompleted status = "NotCompleted"
)

type FileMigration struct {
	Name      string
	AppliedAt *time.Time `db:"applied_at"`
	Status    status
	Content   string
}

type DbMigration struct {
	Name      string
	AppliedAt *time.Time `db:"applied_at"`
	Status    status
}

func Run(pool *pgxpool.Pool) error {
	ctx := &MigratorCtx{pool}
	err := ctx.setup()
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could create migration table")
	}

	migrationFiles, err := getMigrationFiles()
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could not get migration files")
	}

	log.Info().Msgf("Found %d migrations", len(migrationFiles))
	if len(migrationFiles) == 0 {
		return nil
	}

	dbMigrations, err := ctx.getAllFromDb()
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could not fetch migrations from database")
	}

	for _, dbMigration := range dbMigrations {
		if dbMigration.Status == Failure {
			log.Fatal().Msgf(`Found migration with failure status: "%s". Needs manual intervention`, dbMigration.Name)
		}

		migrationFile := migrationFiles[dbMigration.Name]
		if migrationFile != nil {
			migrationFile.Status = dbMigration.Status
			migrationFile.AppliedAt = dbMigration.AppliedAt
		} else {
			log.Fatal().Msgf(`Could not migration file "%s"`, dbMigration.Name)
		}
	}

	for _, migration := range migrationFiles {
		if migration.Status == Success {
			log.Info().Msgf(`Skipping "%s"`, migration.Name)
		} else if migration.Status == Failure {
			log.Fatal().Msgf(`Found failed "%s". Halting.`, migration.Name)
		} else if migration.Status == Pending {
			log.Fatal().Msgf(`Found pending "%s". Indication that previous migration failed, thus halting...`, migration.Name)
		} else if migration.Status == NotCompleted {
			err = ctx.runMigration(migration)
			if err != nil {
				log.Fatal().AnErr("error", err).Msgf(`Migration failed: "%s"`, migration.Name)
			}
		} else {
			log.Fatal().Msgf(`Found unknown migration status: "%s"`, migration.Status)
		}
	}

	return nil
}

func (db MigratorCtx) setup() error {
	ctx, cancel := defaultContext()
	defer cancel()

	var migrationTable = `
		CREATE TABLE IF NOT EXISTS migration
		(
			name       VARCHAR PRIMARY KEY NOT NULL,
			applied_at TIMESTAMP DEFAULT NOW() NOT NULL,
			status     VARCHAR
		);
	`

	_, err := db.pool.Exec(ctx, migrationTable)

	return err
}

func (db MigratorCtx) getAllFromDb() ([]*DbMigration, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	var result []*DbMigration
	err := pgxscan.Select(ctx, db.pool, &result, `SELECT * FROM migration`)

	return result, err
}

func (db MigratorCtx) runMigration(migration *FileMigration) error {
	ctx, cancel := defaultContext()
	defer cancel()

	_, err := db.pool.Exec(ctx, "INSERT INTO migration(name, status) VALUES ($1, $2)", migration.Name, Pending)
	if err != nil {
		log.Fatal().AnErr("error", err).Msgf("Could not create migration entry for %s", migration.Name)
	}

	log.Info().Msgf(`Running migration: %s`, migration.Name)
	_, err = db.pool.Exec(ctx, migration.Content)

	if err != nil {
		log.Fatal().AnErr("error", err).Msgf("Migration failed for %s", migration.Name)
		db.pool.Exec(ctx, "UPDATE migration SET status = $2 WHERE name = $1", migration.Name, Failure)
	} else {
		db.pool.Exec(ctx, "UPDATE migration SET status = $2 WHERE name = $1", migration.Name, Success)
	}

	return err
}

func getMigrationFiles() (map[string]*FileMigration, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(filepath.Join(workingDirectory, "/sql/migrations"))
	if err != nil {
		return nil, err
	}

	out := utils.Reduce(dir, make(map[string]*FileMigration), func(acc map[string]*FileMigration, value os.DirEntry, i int) map[string]*FileMigration {
		content, _ := os.ReadFile(filepath.Join(workingDirectory, "/sql/migrations", value.Name()))
		acc[value.Name()] = &FileMigration{
			Name:      value.Name(),
			AppliedAt: nil,
			Status:    NotCompleted,
			Content:   string(content),
		}
		return acc
	})

	return out, nil
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
