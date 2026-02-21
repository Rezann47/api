package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Rezann47/YksKoc/internal/config"
)

type Runner struct {
	m *migrate.Migrate
}

func New(cfg *config.DBConfig, migrationsPath string) (*Runner, error) {
	m, err := migrate.New(migrationsPath, cfg.MigrateDSN())
	if err != nil {
		return nil, fmt.Errorf("migrate.New: %w", err)
	}
	return &Runner{m: m}, nil
}

func (r *Runner) Up() error {
	if err := r.m.Up(); errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else {
		return err
	}
}

func (r *Runner) Down() error {
	if err := r.m.Down(); errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else {
		return err
	}
}

func (r *Runner) Steps(n int) error            { return r.m.Steps(n) }
func (r *Runner) Force(v int) error            { return r.m.Force(v) }
func (r *Runner) Version() (uint, bool, error) { return r.m.Version() }
func (r *Runner) Close() error {
	s, d := r.m.Close()
	if s != nil {
		return s
	}
	return d
}
