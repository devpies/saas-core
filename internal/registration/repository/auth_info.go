package repository

import (
	"context"

	"github.com/devpies/saas-core/internal/registration/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

// AuthInfoRepository manages data access to tenant authentication information.
type AuthInfoRepository struct {
	client *dynamodb.Client
	logger *zap.Logger
	table  string
}

// NewAuthInfoRepository returns a new AuthInfoRepository.
func NewAuthInfoRepository(logger *zap.Logger, client *dynamodb.Client, table string) *AuthInfoRepository {
	return &AuthInfoRepository{
		logger: logger,
		client: client,
		table:  table,
	}
}

// SelectAuthInfo retrieves authentication information for a specific tenant.
func (ar *AuthInfoRepository) SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error) {
	var authInfo model.AuthInfo
	input := dynamodb.GetItemInput{
		TableName: aws.String(ar.table),
		Key: map[string]types.AttributeValue{
			"tenantPath": &types.AttributeValueMemberS{Value: path},
		},
	}
	output, err := ar.client.GetItem(ctx, &input)
	if err != nil {
		return authInfo, err
	}
	err = attributevalue.UnmarshalMap(output.Item, &authInfo)
	if err != nil {
		return authInfo, err
	}
	return authInfo, nil
}

// InsertAuthInfo stores authentication information required for the tenant login.
func (ar *AuthInfoRepository) InsertAuthInfo(ctx context.Context, info model.AuthInfo) error {
	input := dynamodb.PutItemInput{
		TableName: aws.String(ar.table),
		Item: map[string]types.AttributeValue{
			"tenantPath":       &types.AttributeValueMemberS{Value: info.TenantPath},
			"userPoolId":       &types.AttributeValueMemberS{Value: info.UserPoolID},
			"userPoolType":     &types.AttributeValueMemberS{Value: info.UserPoolType},
			"userPoolClientId": &types.AttributeValueMemberS{Value: info.UserPoolClientID},
		},
	}
	_, err := ar.client.PutItem(ctx, &input)
	if err != nil {
		return err
	}
	return nil
}
