package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Service struct {
	dsn string
}

func New(dsn string) *Service {
	return &Service{
		dsn: dsn,
	}
}

func (s *Service) Up() error {
	m, err := migrate.New("file://migrations", s.dsn)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}

	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("m.Up: %w", err)
		}
	}

	return nil
}

func (s *Service) DownOne() error {
	m, err := migrate.New("file://migrations", s.dsn)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}

	err = m.Steps(-1)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("m.Down: %w", err)
		}
	}

	return nil
}
