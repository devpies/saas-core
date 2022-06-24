package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type registrationClient interface {
	Register(ctx context.Context, tenant []byte) (*http.Response, error)
}

// RegistrationService is responsible for triggering tenant registration.
type RegistrationService struct {
	logger     *zap.Logger
	httpClient registrationClient
}

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger, httpClient registrationClient) *RegistrationService {
	return &RegistrationService{
		logger:     logger,
		httpClient: httpClient,
	}
}

// RegisterTenant sends new tenant to tenant registration microservice and return a nil ErrorResponse on success.
func (rs *RegistrationService) RegisterTenant(ctx context.Context, newTenant model.NewTenant) (*web.ErrorResponse, int, error) {
	var (
		resp *http.Response
		err  error
	)

	tenant := model.Tenant{
		ID:       uuid.New().String(),
		Email:    newTenant.Email,
		FullName: newTenant.FullName,
		Company:  newTenant.Company,
		Plan:     newTenant.Plan,
	}

	data, err := json.Marshal(tenant)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	resp, err = rs.httpClient.Register(ctx, data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var webErrResp web.ErrorResponse
		err = json.Unmarshal(bodyBytes, &webErrResp)
		if err != nil {
			rs.logger.Error("error w/ decoding body", zap.Error(err))
			return nil, resp.StatusCode, err
		}
		return &webErrResp, resp.StatusCode, err
	}

	return nil, resp.StatusCode, nil
}
