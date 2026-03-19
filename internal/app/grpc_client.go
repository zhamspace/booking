package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/zhamspace/booking/internal/config"
	"github.com/zhamspace/booking/internal/errs"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/zhamspace/booking/pkg/proto/common"
)

func newGrpcClientConn(uri string, secure bool, username, password, errMessagePrefix string) (*grpc.ClientConn, error) {
	if uri == "" {
		return nil, nil
	}

	errInterceptor := &grpcClientInterceptorErrorT{
		errMessagePrefix: errMessagePrefix,
	}

	dialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithChainUnaryInterceptor(
			errInterceptor.grpcClientInterceptorError,
		),
	}
	if config.Conf.WithTracing {
		dialOptions = append(dialOptions, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	if secure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if username != "" {
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(
			newGrpcClientBasicAuth(username, password),
		))
	}

	conn, err := grpc.NewClient(uri, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return conn, nil
}

// error interceptor
type grpcClientInterceptorErrorT struct {
	errMessagePrefix string
}

func (o *grpcClientInterceptorErrorT) grpcClientInterceptorError(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		if len(st.Details()) > 0 {
			stDetail := st.Details()[0]
			errObj, ok := stDetail.(*common.ErrorDetail)
			if ok {
				return errs.New(errs.Err(errObj.Code), o.errMessagePrefix+errObj.Message, errObj.Fields)
			}
		}
		return errs.New(errs.ServiceNA, o.errMessagePrefix+st.String(), nil)
	}

	return fmt.Errorf(o.errMessagePrefix+": %w", err)
}

// basic auth

type grpcClientBasicAuth struct {
	metadata map[string]string
}

func newGrpcClientBasicAuth(username, password string) *grpcClientBasicAuth {
	return &grpcClientBasicAuth{
		metadata: map[string]string{
			"authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
		},
	}
}

func (o *grpcClientBasicAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return o.metadata, nil
}

func (o *grpcClientBasicAuth) RequireTransportSecurity() bool {
	return false
}
