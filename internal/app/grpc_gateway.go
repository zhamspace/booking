package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/zhamspace/booking/internal/config"
	"github.com/zhamspace/booking/internal/domain/common/session"
	"github.com/zhamspace/booking/internal/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type GrpcGwServer struct {
	server   *http.Server
	grpcConn *grpc.ClientConn
}

type HttpRoute struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request, map[string]string)
}

func NewGrpcGwServer(
	handlers []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error,
	httpRoutes []HttpRoute,
) (*GrpcGwServer, error) {
	var conn *grpc.ClientConn
	handler, err := GrpcGatewayCreateHandler(func(mux *runtime.ServeMux) error {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		var err error
		conn, err = grpc.NewClient("localhost:"+config.Conf.GrpcPort, opts...)
		if err != nil {
			return fmt.Errorf("grpc.NewClient: %w", err)
		}

		for _, h := range handlers {
			err = h(context.Background(), mux, conn)
			if err != nil {
				return fmt.Errorf("grpc-gateway: register grpc-handler: %w", err)
			}
		}

		for _, r := range httpRoutes {
			if err := mux.HandlePath(r.Method, r.Path, r.Handler); err != nil {
				return fmt.Errorf("mux.HandlePath: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("GrpcGatewayCreateHandler: %w", err)
	}

	// server
	server := &http.Server{
		Addr:              ":" + config.Conf.HttpPort,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       time.Minute,
		MaxHeaderBytes:    300 * 1024,
	}

	return &GrpcGwServer{
		server:   server,
		grpcConn: conn,
	}, nil
}

func (s *GrpcGwServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http-server stopped", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("http-server started " + s.server.Addr)
}

func (s *GrpcGwServer) Stop(ctx context.Context) error {
	ctx, ctxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer ctxCancel()

	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("http-server shutdown error", "error", err)
		return err
	}

	if s.grpcConn != nil {
		if err := s.grpcConn.Close(); err != nil {
			slog.Error("grpc client connection close error", "error", err)
		}
	}

	return nil
}

func GrpcGatewayCreateHandler(muxHook func(*runtime.ServeMux) error) (http.Handler, error) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithErrorHandler(func(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
			repObj, grpcCode := errs.ExtractDetail(err)
			headerStatus := runtime.HTTPStatusFromCode(grpcCode)

			repBody, marshalErr := protojson.Marshal(repObj)
			if marshalErr != nil {
				slog.Error("GRPC_GW: ErrorHandler: Failed to marshal", "error", marshalErr)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(headerStatus)
			if _, err = io.Copy(w, bytes.NewReader(repBody)); err != nil {
				slog.Error("GRPC_GW: ErrorHandler: Failed to write response", "error", err)
			}
		}))

	if muxHook != nil {
		err := muxHook(mux)
		if err != nil {
			return nil, fmt.Errorf("grpc-gateway: muxHook: %w", err)
		}
	}

	// add health check handler
	err := mux.HandlePath(http.MethodGet, "/healthcheck", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-gateway: register healthcheck handler: %w", err)
	}

	// add docs handler
	docFS := http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs")))
	err = mux.HandlePath(http.MethodGet, "/docs/{any}", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		docFS.ServeHTTP(w, r)
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-gateway: register docs handler: %w", err)
	}
	err = mux.HandlePath(http.MethodGet, "/docs/{any}/{any}", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		docFS.ServeHTTP(w, r)
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-gateway: register docs handler: %w", err)
	}

	// add metrics handler
	if err = mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		return nil, fmt.Errorf("grpc-gateway: register metrics handler: %w", err)
	}

	handler := http.Handler(mux)

	// add cors middleware
	if config.Conf.HttpCors {
		handler = cors.New(cors.Options{
			AllowOriginFunc: func(_ string) bool { return true },
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPut,
				http.MethodPost,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"Accept",
				"Content-Type",
				"X-Requested-With",
				"Authorization",
			},
			AllowCredentials: true,
			MaxAge:           604800,
		}).Handler(handler)
	}

	// add recover middleware
	handler = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// use always new err instance in defer
				if err := recover(); err != nil {
					slog.Error(
						"Recovered from panic",
						slog.Any("error", err),
						slog.Any("recovery_stacktrace", string(debug.Stack())),
					)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}(handler)

	// add auth token middleware
	handler = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			if auth := r.Header.Get("Authorization"); len(auth) > 0 {
				token = strings.TrimPrefix(auth, "Bearer ")
			}

			ctx := session.ContextWithToken(r.Context(), token)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}(handler)

	return handler, nil
}
