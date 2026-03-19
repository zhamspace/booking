package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/zhamspace/booking/internal/config"
	appLogging "github.com/zhamspace/booking/internal/domain/common/logging"
	"github.com/zhamspace/booking/internal/domain/common/session"
	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/pkg/proto/common"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type GrpcServer struct {
	name   string
	server *grpc.Server
}

func NewGrpcServer(name string, register func(*grpc.Server)) *GrpcServer {
	interceptors := make([]grpc.UnaryServerInterceptor, 0, 5)
	serverOpts := make([]grpc.ServerOption, 0, 4)

	// map errors to grpc status after logging
	interceptors = append(interceptors, GrpcInterceptorErrToStatus())

	// ctx without cancel
	interceptors = append(interceptors, GrpcInterceptorCtxWithoutCancel())

	// recovery
	interceptors = append(interceptors, GrpcInterceptorRecovery())

	// logging
	interceptors = append(interceptors, GrpcInterceptorLogging())

	// tracing
	if config.Conf.WithTracing {
		serverOpts = append(serverOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	// metrics
	if config.Conf.WithMetrics {
		interceptors = append(interceptors, GrpcInterceptorMetrics())
	}

	// auth
	interceptors = append(interceptors, GrpcInterceptorAuth())

	// server
	serverOpts = append(
		serverOpts,
		grpc.MaxSendMsgSize(math.MaxUint32),
		grpc.MaxRecvMsgSize(math.MaxUint32),
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	server := grpc.NewServer(serverOpts...)

	// register handlers
	if register != nil {
		register(server)
	}

	// register grpc reflection
	reflection.Register(server)

	return &GrpcServer{
		name:   name,
		server: server,
	}
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", ":"+config.Conf.GrpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen grpc: %w", err)
	}
	go func() {
		err = s.server.Serve(lis)
		if err != nil {
			slog.Error(s.name+"-grpc-server stopped", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info(s.name + "-grpc-server started " + lis.Addr().String())
	return nil
}

func (s *GrpcServer) Stop() {
	s.server.GracefulStop()
}

func GrpcInterceptorCtxWithoutCancel() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return handler(context.WithoutCancel(ctx), req)
	}
}

func GrpcInterceptorRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error(
					"Recovered from grpc panic",
					slog.Any("error", recovered),
					slog.String("fullMethod", info.FullMethod),
					slog.Any("recovery_stacktrace", string(debug.Stack())),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

func GrpcInterceptorMetrics() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		h, err := handler(ctx, req)

		responseStatus := "ok"
		if err != nil {
			responseStatus = "error"
		}

		metricRequestCounter.WithLabelValues("grpc", info.FullMethod, responseStatus).Inc()
		metricResponseDuration.WithLabelValues("grpc", info.FullMethod, responseStatus).Observe(time.Since(start).Seconds())

		return h, err
	}
}

func GrpcInterceptorLogging() grpc.UnaryServerInterceptor {
	toJSONString := func(v any) string {
		if v == nil {
			return ""
		}
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("marshal error: %v", err)
		}
		return string(b)
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		clientIP := lo.FirstOr(md.Get("x-forwarded-for"), "unknown")
		ua := lo.FirstOr(md.Get("grpcgateway-user-agent"), lo.FirstOr(md.Get("user-agent"), "unknown"))

		start := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(start)

		mySlog := slog.With(
			slog.String("fullMethod", info.FullMethod),
			slog.String("duration", duration.String()),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", ua),
		)

		policy := appLogging.GetPolicy(info.FullMethod)
		if policy.LogRequest {
			mySlog = mySlog.With(
				"request_body", toJSONString(req))
		}

		if err != nil && policy.LogError {
			ei, stCode := errs.ExtractDetail(err)
			if appErr, ok := errs.AsError(err); ok && appErr.Cause != nil {
				fieldsCopy := map[string]string{}
				for k, v := range ei.Fields {
					fieldsCopy[k] = v
				}
				if _, ok := fieldsCopy["cause"]; !ok {
					fieldsCopy["cause"] = appErr.Cause.Error()
				}
				ei = &common.ErrorDetail{
					Code:    ei.Code,
					Message: ei.Message,
					Fields:  fieldsCopy,
				}
			}
			args := []any{
				slog.String("status", stCode.String()),
				slog.Any("error", ei),
			}
			mySlog.Error("grpc request error", args...)
		}

		if err == nil && policy.LogSuccess {
			if policy.LogResponse {
				mySlog = mySlog.With(
					"response_body", toJSONString(resp))
			}
			mySlog.Info("grpc request ok")
		}

		return resp, err
	}
}

func GrpcInterceptorErrToStatus() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		return resp, errs.ToStatus(err)
	}
}

func GrpcInterceptorAuth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var token string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			auth := md.Get("authorization")
			if len(auth) > 0 {
				token = strings.TrimPrefix(auth[0], "Bearer ")
			}
		}

		ctx = session.ContextWithToken(ctx, token)

		return handler(ctx, req)
	}
}
