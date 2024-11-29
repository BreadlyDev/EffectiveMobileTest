package migrator

import (
	"effectivemobiletesttask/internal/config"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfgDB *config.DBStorage, cfgMigr *config.Migrations) {
	if cfgMigr.Path == "" {
		log.Fatal("migrations-path is required")
	}

	DBUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=%s",
		cfgDB.User, cfgDB.Pass, cfgDB.Host, cfgDB.Port, cfgDB.DBName, cfgDB.SSLMode, cfgMigr.Table)

	m, err := migrate.New(
		"file://"+cfgMigr.Path,
		DBUrl,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Debug("no migrations to apply")
			fmt.Println("no migrations to apply")

			return
		}

		log.Fatal(err)
	}

	log.Println("migrations applied")
}
