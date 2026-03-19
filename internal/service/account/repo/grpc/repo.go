package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/internal/service/account/model"
	"github.com/zhamspace/booking/pkg/proto/account_v1"
	"google.golang.org/grpc"
)

type ConfigSt struct {
	GrpcClient *grpc.ClientConn
}

func (c *ConfigSt) normalize() {
}

func (c *ConfigSt) validate() (finalError error) {
	if c.GrpcClient == nil {
		err := fmt.Errorf("missing field: grpcClient")
		finalError = errors.Join(finalError, err)
	}
	return
}

type Repo struct {
	grpcClient account_v1.AccountClient
	config     *ConfigSt
}

func New(cfg *ConfigSt) (_ *Repo, finalError error) {
	if cfg == nil {
		return nil, errs.InvalidConfig
	}

	cfg.normalize()
	if err := cfg.validate(); err != nil {
		finalError = fmt.Errorf("cfg.validate: %w", err)
	}

	return &Repo{
		grpcClient: account_v1.NewAccountClient(cfg.GrpcClient),
		config:     cfg,
	}, finalError
}

func (r *Repo) GetUser(ctx context.Context, pars *model.GetReq) (*model.Main, bool, error) {
	reqObj := &account_v1.GetUserReq{
		Id: pars.Id,
	}
	repObj, err := r.grpcClient.GetUser(ctx, reqObj)
	if err != nil {
		return nil, false, err
	}

	return &model.Main{
		Id:        repObj.Id,
		Username:  repObj.Username,
		FirstName: repObj.FirstName,
		LastName:  repObj.LastName,
		Email:     repObj.Email,
		CreatedAt: repObj.CreatedAt.AsTime(),
	}, true, nil
}
