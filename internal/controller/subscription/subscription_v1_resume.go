package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/service"
)

// Resume 恢复订阅
func (c *ControllerV1) Resume(ctx context.Context, req *v1.ResumeReq) (*v1.ResumeRes, error) {
	svc := subscription.NewService()
	if err := svc.Resume(ctx, req.ID); err != nil {
		return nil, err
	}

	// 获取订阅信息并添加到调度器
	sub, err := svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	service.AddSubscriptionToScheduler(sub)

	return &v1.ResumeRes{}, nil
}
