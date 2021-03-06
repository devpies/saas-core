package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

// InviteService manages the invite business operations.
type InviteService struct {
	logger     *zap.Logger
	inviteRepo inviteRepository
}

// NewInviteService returns a new invite service.
func NewInviteService(
	logger *zap.Logger,
	inviteRepo inviteRepository,
) *InviteService {
	return &InviteService{
		logger:     logger,
		inviteRepo: inviteRepo,
	}
}

func (is *InviteService) Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error) {
	return is.inviteRepo.Create(ctx, ni, now)
}
func (is *InviteService) RetrieveInvite(ctx context.Context, iid string) (model.Invite, error) {
	return is.inviteRepo.RetrieveInvite(ctx, iid)
}
func (is *InviteService) RetrieveInvites(ctx context.Context) ([]model.Invite, error) {
	return is.inviteRepo.RetrieveInvites(ctx)
}
func (is *InviteService) Update(ctx context.Context, update model.UpdateInvite, iid string, now time.Time) (model.Invite, error) {
	return is.inviteRepo.Update(ctx, update, iid, now)
}
