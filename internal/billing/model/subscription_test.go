package model_test

import (
	"testing"
	"time"

	"github.com/devpies/saas-core/internal/billing/model"

	"github.com/stretchr/testify/assert"
)

func TestNewSubscription_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(ns *model.NewSubscription)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(ns *model.NewSubscription) {},
			err:      "",
		},
		{
			name: "invalid plan",
			modifier: func(ns *model.NewSubscription) {
				ns.Plan = 2
			},
			err: "failed on the 'oneof' tag",
		},
		{
			name: "invalid transaction id",
			modifier: func(ns *model.NewSubscription) {
				ns.TransactionID = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "transaction id is not a uuid",
			modifier: func(ns *model.NewSubscription) {
				ns.TransactionID = model.InvalidUUID
			},
			err: "failed on the 'uuid4' tag",
		},
		{
			name: "invalid status id",
			modifier: func(ns *model.NewSubscription) {
				ns.Plan = 3
			},
			err: "failed on the 'oneof' tag",
		},
		{
			name: "invalid amount",
			modifier: func(ns *model.NewSubscription) {
				ns.Amount = 0
			},
			err: "failed on the 'gt' tag",
		},
		{
			name: "invalid customer id",
			modifier: func(ns *model.NewSubscription) {
				ns.CustomerID = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "customer id is not a uuid",
			modifier: func(ns *model.NewSubscription) {
				ns.CustomerID = model.InvalidUUID
			},
			err: "failed on the 'uuid4' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ns := model.NewSubscription{
				Plan:          1,
				TransactionID: "15ee4790-199a-4514-bb5f-573fe935198e",
				StatusID:      2,
				Amount:        1000,
				CustomerID:    "f0900a2d-5c56-4d77-ba4d-c2165176773c",
			}

			tc.modifier(&ns)

			err := ns.Validate()
			if tc.err != "" {
				if err == nil {
					t.Errorf("expected: %s, got nil", tc.err)
					return
				}
				assert.Regexp(t, tc.err, err.Error())
			} else {
				if err != nil {
					t.Errorf("expected: nil, got: %s", err.Error())
				}
			}
		})
	}
}

func TestUpdateSubscription_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(us *model.UpdateSubscription)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(us *model.UpdateSubscription) {},
			err:      "",
		},
		{
			name: "invalid plan",
			modifier: func(us *model.UpdateSubscription) {
				var plan model.SubscriptionPlanType = 2
				us.Plan = &plan
			},
			err: "failed on the 'oneof' tag",
		},
		{
			name: "invalid transaction id",
			modifier: func(us *model.UpdateSubscription) {
				us.TransactionID = &model.InvalidUUID
			},
			err: "failed on the 'uuid4' tag",
		},
		{
			name: "invalid status id",
			modifier: func(us *model.UpdateSubscription) {
				var id model.SubscriptionStatusType = 3
				us.StatusID = &id
			},
			err: "failed on the 'oneof' tag",
		},
		{
			name: "invalid time",
			modifier: func(us *model.UpdateSubscription) {
				us.UpdatedAt = time.Time{}
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			us := model.UpdateSubscription{
				UpdatedAt: time.Now(),
			}

			tc.modifier(&us)

			err := us.Validate()
			if tc.err != "" {
				if err == nil {
					t.Errorf("expected: %s, got nil", tc.err)
					return
				}
				assert.Regexp(t, tc.err, err.Error())
			} else {
				if err != nil {
					t.Errorf("expected: nil, got: %s", err.Error())
				}
			}
		})
	}
}