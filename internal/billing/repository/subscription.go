// Package repository manages the data access layer for handling queries.
package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/devpies/saas-core/internal/billing/db"
	"github.com/devpies/saas-core/internal/billing/model"
	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// SubscriptionRepository manages data access to subscriptions.
type SubscriptionRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewSubscriptionRepository returns a new SubscriptionRepository.
func NewSubscriptionRepository(logger *zap.Logger, pg *db.PostgresDatabase) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger: logger,
		pg:     pg,
	}
}

// SaveSubscription saves a subscription.
func (sr *SubscriptionRepository) SaveSubscription(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error) {
	var (
		s   model.Subscription
		err error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return s, web.CtxErr()
	}
	conn, Close, err := sr.pg.GetConnection(ctx)
	if err != nil {
		return s, err
	}
	defer Close()

	stmt := `
			insert into subscriptions (
   				subscription_id, plan, transaction_id, subscription_status_id,
				amount, customer_id, tenant_id
			) values ($1, $2, $3, $4, $5, $6, $7)
		`
	s = model.Subscription{
		ID:            uuid.New().String(),
		Plan:          ns.Plan,
		TransactionID: ns.TransactionID,
		StatusID:      ns.StatusID,
		Amount:        ns.Amount,
		TenantID:      values.TenantID,
		CustomerID:    ns.CustomerID,
		UpdatedAt:     now.UTC(),
		CreatedAt:     now.UTC(),
	}

	if _, err = conn.ExecContext(
		ctx,
		stmt,
		s.ID,
		s.Plan,
		s.TransactionID,
		s.StatusID,
		s.Amount,
		s.CustomerID,
		s.TenantID,
	); err != nil {
		return s, fmt.Errorf("error inserting subscription %v :%w", s, err)
	}

	return s, nil
}

// GetAllSubscriptions retrieves all subscriptions.
func (sr *SubscriptionRepository) GetAllSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	var subs []model.Subscription
	return subs, nil
}

// GetOneSubscription retrieves one subscription by id.
func (sr *SubscriptionRepository) GetOneSubscription(ctx context.Context, id string) (model.Subscription, error) {
	var s model.Subscription
	return s, nil
}