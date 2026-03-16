package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/service"
)

// Create 创建订阅
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	svc := subscription.NewService()

	sub := &model.Subscription{
		Name:        req.Name,
		Description: req.Description,
		SourceType:  model.SourceType(req.SourceType),
		SourceURL:   req.SourceURL,
		CronExpr:    req.CronExpr,
		Config:      req.Config,
		Status:      model.StatusActive,
	}

	if err := svc.Create(ctx, sub); err != nil {
		return nil, err
	}

	// 添加到调度器
	service.AddSubscriptionToScheduler(sub)

	return &v1.CreateRes{ID: sub.ID}, nil
}
