package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type userRepository interface {
	Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error)
	CreateAdmin(ctx context.Context, na model.NewAdminUser) error
	List(ctx context.Context) ([]model.User, error)
	RetrieveByEmail(ctx context.Context, email string) (model.User, error)
	RetrieveMe(ctx context.Context) (model.User, error)
}

type cognitoClient interface {
	AdminCreateUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminCreateUserInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

// UserService manages the user business operations.
type UserService struct {
	logger        *zap.Logger
	userPoolID    string
	userRepo      userRepository
	cognitoClient cognitoClient
}

// NewUserService returns a new user service.
func NewUserService(
	logger *zap.Logger,
	userPoolID string,
	userRepo userRepository,
	cognitoClient cognitoClient,
) *UserService {
	return &UserService{
		logger:        logger,
		userPoolID:    userPoolID,
		userRepo:      userRepo,
		cognitoClient: cognitoClient,
	}
}

// AddSeat publishes a message to create a user in the identity service.
func (us *UserService) AddSeat(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error) {
	var (
		u   model.User
		err error
	)
	_, err = us.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(us.userPoolID),
		Username:   aws.String(nu.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(uuid.New().String())},
			{Name: aws.String("custom:company-name"), Value: aws.String(nu.Company)},
			{Name: aws.String("custom:full-name"), Value: aws.String(fmt.Sprintf("%s %s", nu.FirstName, nu.LastName))},
			{Name: aws.String("email"), Value: aws.String(nu.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
	if err != nil {
		us.logger.Error("", zap.Error(err))
		return u, err
	}
	us.logger.Info("successfully added user")

	user, err := us.userRepo.Create(ctx, nu, now)
	if err != nil {
		return u, err
	}
	return user, nil
}

func (us *UserService) AddAdminUserFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalTenantIdentityCreatedEvent(m)
	if err != nil {
		return err
	}
	na := newAdminUser(event.Data)
	return us.userRepo.CreateAdmin(ctx, na)
}

func newAdminUser(data msg.TenantIdentityCreatedEventData) model.NewAdminUser {
	return model.NewAdminUser{
		UserID:        data.UserID,
		TenantID:      data.TenantID,
		Company:       data.Company,
		Email:         data.Email,
		FirstName:     data.FirstName,
		LastName:      data.LastName,
		EmailVerified: true,
		CreatedAt:     msg.ParseTime(data.CreatedAt),
	}
}

func (us *UserService) List(ctx context.Context) ([]model.User, error) {
	return us.userRepo.List(ctx)
}

func (us *UserService) RetrieveByEmail(ctx context.Context, email string) (model.User, error) {
	return us.userRepo.RetrieveByEmail(ctx, email)
}

func (us *UserService) RetrieveMe(ctx context.Context) (model.User, error) {
	return us.userRepo.RetrieveMe(ctx)
}
