package command

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pgCon *pgxpool.Pool
}

func New(
	pgCon *pgxpool.Pool,
) *Service {
	return &Service{
		pgCon: pgCon,
	}
}

func recoverFn() {
	if recovered := recover(); recovered != nil {
		slog.Error(
			"Recovered from cmd panic",
			slog.Any("error", recovered),
			slog.Any("recovery_stacktrace", string(debug.Stack())),
		)
	}
}

func (s *Service) Run(num int) error {

	switch num {
	case 1:
		go s.process1()
	case 2:
		go s.process2()
	case 3:
		go s.process3()
	case 4:
		go s.process4()
	default:
		return fmt.Errorf("no process: %d", num)
	}

	return nil
}

func (s *Service) process1() {
	defer recoverFn()
	const op = "cmd1.Name"

	startTime := time.Now()

	slog.Info(op+": start", "id", "1")
	defer func() {
		slog.Info(op+": done", "duration", time.Since(startTime).String(), "id", "1")
	}()

	ctx := context.Background()
	_ = ctx
}

func (s *Service) process2() {
	defer recoverFn()
	const op = "cmd2.Name"
	startTime := time.Now()

	slog.Info(op+": start", "id", "2")
	defer func() {
		slog.Info(op+": done", "duration", time.Since(startTime).String(), "id", "2")
	}()

	ctx := context.Background()
	_ = ctx

	panic(op + ": Bada-boom")
}

func (s *Service) process3() {
	defer recoverFn()
	const op = "cmd3.Name"
	startTime := time.Now()

	slog.Info(op+": start", "id", "3")
	defer func() {
		slog.Info("command 3 done", "duration", time.Since(startTime).String(), "id", "3")
	}()

	ctx := context.Background()
	_ = ctx
}

func (s *Service) process4() {
	defer recoverFn()
	const op = "cmd4.Name"
	startTime := time.Now()

	slog.Info(op+": start", "id", "4")
	defer func() {
		slog.Info("command 4 done", "duration", time.Since(startTime).String(), "id", "4")
	}()

	ctx := context.Background()
	_ = ctx
}
