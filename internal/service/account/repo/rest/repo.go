package rest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/internal/service/account/model"
	repoModel "github.com/zhamspace/booking/internal/service/account/repo/rest/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	GetUserPathTmpl = "users/%s"
)

type ConfigSt struct {
	Uri       string
	SecretKey string
}

func (c *ConfigSt) normalize() {
	c.Uri = fmt.Sprintf("%s/", strings.TrimRight(c.Uri, "/"))
}

func (c *ConfigSt) validate() (finalError error) {
	if c.Uri == "" || c.Uri == "/" {
		err := fmt.Errorf("missing field: Uri")
		finalError = errors.Join(finalError, err)
	}
	return
}

type Repo struct {
	httpClient *http.Client
	config     *ConfigSt
	tracer     trace.Tracer
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
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				TLSHandshakeTimeout: 10 * time.Second,
				MaxIdleConnsPerHost: 100,
			},
		},
		config: cfg,
		tracer: otel.Tracer("venue/account-repo"),
	}, finalError
}

func (r *Repo) GetUser(ctx context.Context, pars *model.GetReq) (*model.Main, bool, error) {
	const op = "account.Repo.GetUser"

	repObj := &repoModel.GetUserRepSt{}
	statusOk, statusCode, respBody, err := r.sendRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(GetUserPathTmpl, pars.Id),
		nil,
		url.Values{},
		nil,
		repObj,
		nil,
	)
	if err != nil {
		return nil, false, err
	}

	if !statusOk {
		return nil, false, fmt.Errorf(op+": !statusOk statusCode=%d resp_body=%s", statusCode, string(respBody))
	}

	return &model.Main{
		Id:        repObj.Id,
		Username:  repObj.Username,
		FirstName: repObj.FirstName,
		LastName:  repObj.LastName,
		Email:     repObj.Email,
		CreatedAt: repObj.CreatedAt,
	}, true, nil
}

func (r *Repo) sendRequest(
	ctx context.Context,
	method string,
	endpoint string,
	headers http.Header,
	params url.Values,
	reqObj any,
	repObj any,
	errRepObj any,
) (statusOk bool, statusCode int, _ []byte, finalError error) {
	uri := fmt.Sprintf("%s%s", r.config.Uri, endpoint)
	ctx, span := r.tracer.Start(ctx, "account: ("+method+") "+uri)
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	var reqStream io.Reader
	if reqObj != nil {
		reqJSON, err := json.Marshal(reqObj)
		if err != nil {
			return false, 0, nil, fmt.Errorf("json.Marshal method=%s uri=%s: %w", method, uri, err)
		}
		reqStream = bytes.NewReader(reqJSON)
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, reqStream)
	if err != nil {
		return false, 0, nil, fmt.Errorf("http.NewRequest method=%s uri=%s: %w", method, uri, err)
	}

	if headers != nil {
		req.Header = headers.Clone()
	}
	if reqObj != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if r.config.SecretKey != "" {
		req.Header.Set("Authorization", fmt.Sprint("Bearer ", r.config.SecretKey))
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return false, 0, nil, fmt.Errorf("http.Client Do method=%s uri=%s: %w", method, uri, err)
	}
	defer func() { _ = resp.Body.Close() }()

	repBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, nil, fmt.Errorf("io.ReadAll method=%s uri=%s: %w", method, uri, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if errRepObj != nil && len(repBody) > 0 {
			_ = json.Unmarshal(repBody, errRepObj)
		}
		return false, resp.StatusCode, repBody, nil
	}

	if repObj != nil {
		err = json.Unmarshal(repBody, repObj)
		if err != nil {
			return false, 0, nil, fmt.Errorf("json.Unmarshal method=%s uri=%s body=%s: %w", method, uri, string(repBody), err)
		}
	}

	return true, resp.StatusCode, repBody, nil
}
