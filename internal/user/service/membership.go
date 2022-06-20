package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"time"
)

type membershipRepository interface {
	Create(ctx context.Context, nm model.NewMembership, now time.Time) (model.Membership, error)
	RetrieveMemberships(ctx context.Context, tid string) ([]model.MembershipEnhanced, error)
	RetrieveMembership(ctx context.Context, tid string) (model.Membership, error)
	Update(ctx context.Context, tid string, update model.UpdateMembership, now time.Time) error
	Delete(ctx context.Context, tid string) (string, error)
}

// MembershipService manages the membership business operations.
type MembershipService struct {
	logger         *zap.Logger
	membershipRepo membershipRepository
}

// NewMembershipService returns a new membership service.
func NewMembershipService(
	logger *zap.Logger,
	membershipRepo membershipRepository,
) *MembershipService {
	return &MembershipService{
		logger:         logger,
		membershipRepo: membershipRepo,
	}
}

func (ms *MembershipService) Create(ctx context.Context, nm model.NewMembership, now time.Time) (model.Membership, error) {
	return ms.membershipRepo.Create(ctx, nm, now)
}

func (ms *MembershipService) RetrieveMemberships(ctx context.Context, tid string) ([]model.MembershipEnhanced, error) {
	return ms.membershipRepo.RetrieveMemberships(ctx, tid)
}

func (ms *MembershipService) RetrieveMembership(ctx context.Context, tid string) (model.Membership, error) {
	return ms.membershipRepo.RetrieveMembership(ctx, tid)
}

func (ms *MembershipService) Update(ctx context.Context, tid string, update model.UpdateMembership, now time.Time) error {
	return ms.membershipRepo.Update(ctx, tid, update, now)
}

func (ms *MembershipService) Delete(ctx context.Context, tid string) (string, error) {
	return ms.membershipRepo.Delete(ctx, tid)
}