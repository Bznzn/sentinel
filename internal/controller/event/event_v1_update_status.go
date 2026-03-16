package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/repository"
)

// UpdateStatus 更新事件状态
func (c *Controller) UpdateStatus(ctx context.Context, req *v1.UpdateStatusReq) (*v1.UpdateStatusRes, error) {
	repo := repository.NewEventRepository()

	// 如果状态为已解决且提供了解决方案，更新 recommendation 字段
	if req.Status == "resolved" && req.Solution != "" {
		event, err := c.service.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		event.Recommendation = req.Solution
		event.Status = model.EventStatus(req.Status)
		if err := repo.Update(event); err != nil {
			return nil, err
		}
	} else {
		if err := repo.UpdateStatus(req.ID, model.EventStatus(req.Status)); err != nil {
			return nil, err
		}
	}

	return &v1.UpdateStatusRes{}, nil
}
