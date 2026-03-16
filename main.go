package main

import (
	_ "SuperBizAgent/internal/ai/skills/builtin" // 空导入，注册内置技能
	"SuperBizAgent/internal/controller/chat"
	"SuperBizAgent/internal/controller/event"
	"SuperBizAgent/internal/controller/report"
	"SuperBizAgent/internal/controller/skill"
	"SuperBizAgent/internal/controller/subscription"
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/service"
	"SuperBizAgent/utility/common"
	"SuperBizAgent/utility/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	// 1. 上下文初始化
	ctx := gctx.New()

	// 2. 配置文件读取
	fileDir, err := g.Cfg().Get(ctx, "file_dir")
	if err != nil {
		panic(err)
	}
	common.FileDir = fileDir.String()

	// 3. 数据库初始化
	if err := database.Init(); err != nil {
		panic(err)
	}

	// 4. 启动后台服务
	service.InitFetcher()

	// 5. 获取 HTTP 服务器实例
	s := g.Server()

	// 6. 注册 SSE 路由（不经过中间件）
	s.BindHandler("POST:/api/event/pipeline/stream", event.New().PipelineStream)

	// 7. 路由分组与中间件
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORSMiddleware)
		group.Middleware(middleware.ResponseMiddleware)
		group.Bind(chat.NewV1())                   //聊天接口（/api/chat/v1/*）
		group.Bind(subscription.NewV1())           //订阅管理（/api/subscription/v1/*）
		group.Bind(report.New())                   //报告管理（/api/report/*）
		group.Bind(report.NewTemplateController()) //报告模板（/api/report/template/*
		group.Bind(event.New())                    //事件管理（/api/event/*）
		group.Bind(skill.NewV1())                  //技能执行（/api/skill/v1/*）
	})

	// 8. 设置端口并启动
	s.SetPort(6872)
	s.Run()
}
