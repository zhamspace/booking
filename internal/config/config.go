package config

import (
	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

var Conf = struct {
	Namespace    string `env:"NAMESPACE" envDefault:"example.com"`
	Debug        bool   `env:"DEBUG" envDefault:"false"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
	GrpcPort     string `env:"GRPC_PORT" envDefault:"5050"`
	HttpPort     string `env:"HTTP_PORT" envDefault:"80"`
	HttpCors     bool   `env:"HTTP_CORS" envDefault:"false"`
	TestMode     bool   `env:"TEST_MODE" envDefault:"false"`
	WithMetrics  bool   `env:"WITH_METRICS" envDefault:"false"`
	WithTracing  bool   `env:"WITH_TRACING" envDefault:"false"`
	OtlpEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	PgDsn        string `env:"PG_DSN"`

	AccountGrpcUrl      string `env:"ACCOUNT_GRPC_URL"`
	AccountGrpcSecure   bool   `env:"ACCOUNT_GRPC_SECURE" envDefault:"false"`
	AccountGrpcUsername string `env:"ACCOUNT_GRPC_USERNAME"`
	AccountGrpcPassword string `env:"ACCOUNT_GRPC_PASSWORD"`
	AccountHttpUrl      string `env:"ACCOUNT_HTTP_URL"`
	SessionHttpUrl      string `env:"SESSION_HTTP_URL" envDefault:"http://session:80"`
}{}

func init() {
	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
