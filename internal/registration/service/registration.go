package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/registration/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"go.uber.org/zap"
)

// RegistrationService is responsible for managing tenant registration.
type RegistrationService struct {
	logger     *zap.Logger
	region     string
	idpService identityService
	js         publisher
}

type identityService interface {
	GetPlanBasedUserPool(ctx context.Context, tenant model.NewTenant, path string) (string, error)
}

// Plan represents the type of subscription plan.
type Plan string

const (
	// PlanBasic represents the cheapest subscription plan offering.
	PlanBasic Plan = "basic"
	// PlanPremium represents the premium plan offering.
	PlanPremium Plan = "premium"
)

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger, region string, idpService identityService, js publisher) *RegistrationService {
	return &RegistrationService{
		logger:     logger,
		region:     region,
		idpService: idpService,
		js:         js,
	}
}

// CreateRegistration starts the tenant registration process.
func (rs *RegistrationService) CreateRegistration(ctx context.Context, id string, tenant model.NewTenant) error {
	var err error
	userPoolID, err := rs.idpService.GetPlanBasedUserPool(ctx, tenant, formatPath(tenant.Company))
	if err != nil {
		return err
	}
	err = rs.publishTenantRegisteredEvent(ctx, id, tenant, userPoolID)
	if err != nil {
		return err
	}
	if err = rs.provision(ctx, Plan(tenant.Plan)); err != nil {
		return err
	}
	return nil
}

func formatPath(company string) string {
	return strings.ToLower(strings.Replace(company, " ", "", -1))
}

func (rs *RegistrationService) publishTenantRegisteredEvent(ctx context.Context, id string, tenant model.NewTenant, userPoolID string) error {
	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}
	event := newTenantRegisteredEvent(values, id, tenant, userPoolID)
	bytes, err := event.Marshal()
	if err != nil {
		return err
	}
	rs.js.Publish(msg.SubjectRegistered, bytes)
	return nil
}

func (rs *RegistrationService) provision(ctx context.Context, plan Plan) error {
	if plan != PlanPremium {
		return nil
	}
	client := codepipeline.New(codepipeline.Options{
		Region: rs.region,
	})

	input := codepipeline.StartPipelineExecutionInput{
		Name:               aws.String("tenant-onboarding-pipeline"),
		ClientRequestToken: aws.String(fmt.Sprintf("requestToken-%s", time.Now().UTC())),
	}

	output, err := client.StartPipelineExecution(ctx, &input)
	if err != nil {
		return err
	}
	rs.logger.Info(fmt.Sprintf("successfully started pipeline - response: %+v", output))
	return nil
}

func newTenantRegisteredEvent(values *web.Values, id string, tenant model.NewTenant, userPoolID string) msg.TenantRegisteredEvent {
	return msg.TenantRegisteredEvent{
		Metadata: msg.Metadata{
			TraceID: values.TraceID,
			UserID:  values.UserID,
		},
		Type: msg.TypeTenantRegistered,
		Data: msg.TenantRegisteredEventData{
			ID:         id,
			FullName:   tenant.FullName,
			Company:    tenant.Company,
			Email:      tenant.Email,
			Plan:       tenant.Plan,
			UserPoolID: userPoolID,
		},
	}
}
