package service

import (
	"context"
	"log"
	"time"

	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/fetcher"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/scheduler"
)

var (
	fetchManager    *fetcher.Manager
	globalScheduler *scheduler.Scheduler
)

func InitFetcher() {
	fetchManager = fetcher.NewManager()
	fetchManager.Register(fetcher.NewRSSFetcher())
	fetchManager.Register(fetcher.NewGitHubFetcher())
	fetchManager.Register(fetcher.NewNVDFetcher(""))

	// 初始化调度器
	globalScheduler = scheduler.New(schedulerHandler)
	globalScheduler.Start()

	// 加载所有激活的订阅
	loadActiveSubscriptions()
}

// schedulerHandler 调度器处理函数
func schedulerHandler(ctx context.Context, sub *model.Subscription) error {
	_, err := FetchSubscription(sub)
	return err
}

// loadActiveSubscriptions 加载所有激活的订阅到调度器
func loadActiveSubscriptions() {
	ctx := context.Background()
	svc := subscription.NewService()

	subs, err := svc.GetActiveSubscriptions(ctx)
	if err != nil {
		log.Printf("failed to load active subscriptions: %v", err)
		return
	}

	for _, sub := range subs {
		globalScheduler.Add(sub)
	}

	log.Printf("loaded %d active subscriptions to scheduler", len(subs))
}

// AddSubscriptionToScheduler 添加订阅到调度器
func AddSubscriptionToScheduler(sub *model.Subscription) {
	if globalScheduler != nil {
		globalScheduler.Add(sub)
	}
}

// RemoveSubscriptionFromScheduler 从调度器移除订阅
func RemoveSubscriptionFromScheduler(id uint) {
	if globalScheduler != nil {
		globalScheduler.Remove(id)
	}
}

func runFetchLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fetchAllActive()
	}
}

func fetchAllActive() {
	db := database.GetDB()
	var subs []model.Subscription
	db.Where("status = ?", model.StatusActive).Find(&subs)

	// Worker pool限制并发数
	sem := make(chan struct{}, 10)
	for _, sub := range subs {
		sem <- struct{}{}
		go func(s model.Subscription) {
			defer func() { <-sem }()
			if _, err := FetchSubscription(&s); err != nil {
				log.Printf("fetch subscription err: %v", err)
			}
		}(sub)
	}
}

// FetchResult 抓取结果
type FetchResult struct {
	FetchedCount int   // 抓取到的数量
	NewCount     int   // 新增数量
	TotalEvents  int64 // 总事件数
	Duration     int   // 耗时ms
}

func FetchSubscription(sub *model.Subscription) (*FetchResult, error) {
	if fetchManager == nil {
		log.Println("fetchManager is not initialized")
		return &FetchResult{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	events, err := fetchManager.Fetch(ctx, sub)
	duration := int(time.Since(start).Milliseconds())

	db := database.GetDB()
	logg := &model.FetchLog{
		SubscriptionID: sub.ID,
		Duration:       duration,
		CreatedAt:      time.Now(),
	}

	if err != nil {
		logg.Status = model.FetchStatusFailed
		logg.ErrorMsg = err.Error()
		db.Create(logg)
		return &FetchResult{Duration: duration}, err
	}

	logg.Status = model.FetchStatusSuccess
	logg.EventCount = len(events)
	db.Create(logg)

	// 保存事件(去重)
	newCount := 0
	for _, e := range events {
		var count int64
		db.Model(&model.SecurityEvent{}).Where("unique_hash = ?", e.UniqueHash).Count(&count)
		if count == 0 {
			db.Create(&e)
			newCount++
		}
	}

	// 更新订阅统计 - 累计总事件数
	var totalEvents int64
	db.Model(&model.SecurityEvent{}).Where("subscription_id = ?", sub.ID).Count(&totalEvents)
	now := time.Now()
	db.Model(sub).Updates(map[string]interface{}{
		"last_fetch_at": now,
		"total_events":  totalEvents,
	})

	return &FetchResult{
		FetchedCount: len(events),
		NewCount:     newCount,
		TotalEvents:  totalEvents,
		Duration:     duration,
	}, nil
}
