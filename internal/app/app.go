package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhamspace/booking/internal/config"
	"github.com/zhamspace/booking/internal/constant"

	//domain
	domainAuditServiceP "github.com/zhamspace/booking/internal/domain/audit"
	domainAuditRepoDataP "github.com/zhamspace/booking/internal/domain/audit/repo/db"
	appLogging "github.com/zhamspace/booking/internal/domain/common/logging"

	domainBookingServiceP "github.com/zhamspace/booking/internal/domain/booking"
	domainBookingRepoDataP "github.com/zhamspace/booking/internal/domain/booking/repo/db"
	domainBookingRepoMockP "github.com/zhamspace/booking/internal/domain/booking/repo/mock"

	domainDictServiceP "github.com/zhamspace/booking/internal/domain/common/dict"
	domainDictRepoDataP "github.com/zhamspace/booking/internal/domain/common/dict/repo/data"

	// service
	serviceAcocuntP "github.com/zhamspace/booking/internal/service/account"
	serviceAuditCodecP "github.com/zhamspace/booking/internal/service/auditcodec"
	serviceSessionP "github.com/zhamspace/booking/internal/service/session"
	serviceCommandP "github.com/zhamspace/booking/internal/service/system/command"
	serviceMigrationP "github.com/zhamspace/booking/internal/service/system/migration"

	// usecase
	usecaseAuditP "github.com/zhamspace/booking/internal/usecase/audit"
	usecaseBookingP "github.com/zhamspace/booking/internal/usecase/booking"
	usecaseDictP "github.com/zhamspace/booking/internal/usecase/dict"
	usecaseSystemP "github.com/zhamspace/booking/internal/usecase/system"

	handlerGrpcP "github.com/zhamspace/booking/internal/handler/grpc"
	handlerHttpP "github.com/zhamspace/booking/internal/handler/http"
	expiryWorkerP "github.com/zhamspace/booking/internal/worker/expiry"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"github.com/zhamspace/booking/pkg/proto/booking_v1"
)

type App struct {
	globalTracerCloser io.Closer

	pgpool *pgxpool.Pool

	grpcServer   *GrpcServer
	httpServer   *GrpcGwServer
	auditService *domainAuditServiceP.Service
	expiryWorker *expiryWorkerP.Worker

	ctx       context.Context
	ctxCancel context.CancelFunc

	exitCode int
}

func (a *App) Init() {
	var err error

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	// domain
	var (
		domainAuditService   *domainAuditServiceP.Service
		domainBookingService *domainBookingServiceP.Service
		domainDictService    *domainDictServiceP.Service
	)

	// service
	var (
		accountService *serviceAcocuntP.Service
	)

	// handlers
	var (
		handlerAudit   *handlerGrpcP.Audit
		handlerBooking *handlerGrpcP.Booking
		handlerDict    *handlerGrpcP.Dict

		handlerHttpSystem       *handlerHttpP.System
		handlerHttpInternal     *handlerHttpP.Internal
		handlerHttpReservations *handlerHttpP.Reservations
	)

	var bookingUsecase *usecaseBookingP.Usecase

	// logger
	{
		if err := appLogging.InitLevel(config.Conf.LogLevel); err != nil {
			slog.Error("invalid log level, fallback to info", "error", err)
		}
		slogOpts := &slog.HandlerOptions{
			AddSource: config.Conf.Debug,
			Level:     appLogging.LevelVar(),
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
		if config.Conf.Debug {
			logger = slog.New(slog.NewTextHandler(os.Stdout, slogOpts))
		}
		slog.SetDefault(logger)
	}

	// globalTracer
	{
		if config.Conf.WithTracing && config.Conf.OtlpEndpoint != "" {
			slog.Info("tracing enabled")
			a.globalTracerCloser, err = tracerInitGlobal(config.Conf.OtlpEndpoint, constant.ServiceName)
			errCheck(err, "tracerInitGlobal")
		}
	}

	// pgpool
	{
		pgConf, err := pgxpool.ParseConfig(config.Conf.PgDsn)
		errCheck(err, "pgxpool.ParseConfig")

		pgConf.MaxConns = 10
		pgConf.MinConns = 2
		pgConf.MaxConnLifetime = 3 * time.Minute
		pgConf.MaxConnIdleTime = time.Minute
		pgConf.HealthCheckPeriod = 15 * time.Second

		a.pgpool, err = pgxpool.NewWithConfig(context.Background(), pgConf)
		errCheck(err, "pgxpool.NewWithConfig")
	}

	// dict
	{
		repoData := domainDictRepoDataP.New()
		domainDictService = domainDictServiceP.New(repoData)
		usecase := usecaseDictP.New(domainDictService)
		handlerDict = handlerGrpcP.NewDict(usecase)
	}

	// account
	{
		var repo serviceAcocuntP.RepoI
		repo, err = NewAccountRepo()
		errCheck(err, "NewAccountRepo")

		accountService = serviceAcocuntP.New(repo)
	}

	// audit
	{
		repoData := domainAuditRepoDataP.New(a.pgpool)
		serializer := serviceAuditCodecP.JSONSerializer{}
		domainAuditService = domainAuditServiceP.New(repoData, serializer)
		a.auditService = domainAuditService
		usecase := usecaseAuditP.New(
			accountService,
			domainAuditService,
		)
		handlerAudit = handlerGrpcP.NewAudit(usecase, serializer)
	}

	// booking
	{
		repoData := domainBookingRepoDataP.New(a.pgpool)
		if config.Conf.TestMode {
			domainBookingService = domainBookingServiceP.New(domainBookingRepoMockP.New())
		} else {
			domainBookingService = domainBookingServiceP.New(repoData)
		}
		bookingUsecase = usecaseBookingP.New(domainAuditService, domainBookingService)
		handlerBooking = handlerGrpcP.NewBooking(bookingUsecase)
		handlerHttpInternal = handlerHttpP.NewInternal(bookingUsecase)
		handlerHttpReservations = handlerHttpP.NewReservations(bookingUsecase, serviceSessionP.New(config.Conf.SessionHttpUrl))
		a.expiryWorker = expiryWorkerP.New(domainBookingService, 60*time.Second)
	}

	// http handler
	{
		// system
		{
			migrationService := serviceMigrationP.New(config.Conf.PgDsn)
			loggingService := appLogging.NewManager()
			commandService := serviceCommandP.New(
				a.pgpool,
			)
			usecase := usecaseSystemP.New(commandService, migrationService, loggingService)
			handlerHttpSystem = handlerHttpP.NewSystem(usecase)
		}
	}

	// grpc server
	{
		a.grpcServer = NewGrpcServer("main", func(server *grpc.Server) {
			// server registers
			booking_v1.RegisterAuditServer(server, handlerAudit)
			booking_v1.RegisterBookingServer(server, handlerBooking)
			booking_v1.RegisterDictServer(server, handlerDict)
		})
	}

	// http-gw server
	{
		a.httpServer, err = NewGrpcGwServer(
			[]func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error{
				// server registers
				booking_v1.RegisterAuditHandler,
				booking_v1.RegisterBookingHandler,
				booking_v1.RegisterDictHandler,
			},
			[]HttpRoute{
				{Method: http.MethodGet, Path: "/system/ping", Handler: func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
					w.WriteHeader(http.StatusOK)
					_, _ = io.Copy(w, bytes.NewReader([]byte("pong")))
				}},
				{Method: http.MethodGet, Path: "/system/migration/up", Handler: handlerHttpSystem.MigrationUp},
				{Method: http.MethodGet, Path: "/system/migration/down/one", Handler: handlerHttpSystem.MigrationDownOne},
				{Method: http.MethodGet, Path: "/system/cmd/{num}", Handler: handlerHttpSystem.Command},
				{Method: http.MethodPost, Path: "/system/log/level", Handler: handlerHttpSystem.SetLogLevel},
				{Method: http.MethodPost, Path: "/system/log/policy", Handler: handlerHttpSystem.UpdateLogPolicy},
				{Method: http.MethodDelete, Path: "/system/log/policy/{method}", Handler: handlerHttpSystem.DeleteLogPolicy},
				{Method: http.MethodPost, Path: "/system/log/default", Handler: handlerHttpSystem.UpdateDefaultLogPolicy},
				{Method: http.MethodGet, Path: "/reservations", Handler: handlerHttpReservations.List},
				// Internal routes — called by payment service on internal Docker network, no JWT
				{Method: http.MethodGet, Path: "/internal/booking-stats", Handler: handlerHttpInternal.Stats},
				{Method: http.MethodGet, Path: "/internal/bookings-occupied", Handler: handlerHttpInternal.ListOccupied},
				{Method: http.MethodGet, Path: "/internal/bookings/{id}", Handler: handlerHttpInternal.GetBooking},
				{Method: http.MethodPost, Path: "/internal/bookings/{id}/mark-paid", Handler: handlerHttpInternal.MarkPaid},
				{Method: http.MethodPost, Path: "/internal/bookings/{id}/mark-failed", Handler: handlerHttpInternal.MarkFailed},
				{Method: http.MethodPost, Path: "/internal/bookings/{id}/mark-refunded", Handler: handlerHttpInternal.MarkRefunded},
			})
		errCheck(err, "NewGrpcGwServer")
	}
}

func (a *App) PreStartHook() {
	slog.Info("PreStartHook")
}

func (a *App) Start() {
	slog.Info("Starting", "service", constant.ServiceName)

	if a.auditService != nil {
		if err := a.auditService.Start(context.Background()); err != nil {
			slog.Error("auditService.Start", "error", err)
			a.exitCode = 1
		}
	}

	if a.expiryWorker != nil {
		go a.expiryWorker.Run(a.ctx)
	}

	// grpc server
	{
		err := a.grpcServer.Start()
		errCheck(err, "grpcServer.Start")
	}

	// grpc-gw server
	{
		a.httpServer.Start()
	}
}

func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// wait signal
	sig := <-signalCtx.Done()
	slog.Info("Shutdown signal received", "signal", sig)
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	// stop context
	a.ctxCancel()

	// grpc-gw server
	{
		err := a.httpServer.Stop(context.Background())
		if err != nil {
			slog.Error("httpServer.Stop", "error", err)
			a.exitCode = 1
		}
	}

	// grpc server
	a.grpcServer.Stop()

	if a.auditService != nil {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer stopCancel()
		if err := a.auditService.Stop(stopCtx); err != nil {
			slog.Error("auditService.Stop", "error", err)
			a.exitCode = 1
		}
	}

	if a.pgpool != nil {
		a.pgpool.Close()
	}
}

func (a *App) WaitJobs() {
	slog.Info("waiting jobs")

	if a.auditService != nil {
		waitCtx, waitCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer waitCancel()
		if err := a.auditService.Wait(waitCtx); err != nil {
			slog.Error("auditService.Wait", "error", err)
			a.exitCode = 1
		}
	}
}

func (a *App) Exit() {
	slog.Info("Exit")

	if a.globalTracerCloser != nil {
		_ = a.globalTracerCloser.Close()
	}

	// flush stdout

	os.Exit(a.exitCode)
}

func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
