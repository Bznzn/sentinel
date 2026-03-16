package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/service"
)

// Pause 暂停订阅
func (c *ControllerV1) Pause(ctx context.Context, req *v1.PauseReq) (*v1.PauseRes, error) {
	svc := subscription.NewService()
	if err := svc.Pause(ctx, req.ID); err != nil {
		return nil, err
	}

	// 从调度器移除
	service.RemoveSubscriptionFromScheduler(req.ID)

	return &v1.PauseRes{}, nil
}
