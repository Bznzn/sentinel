# main.go 启动机制详解

## 文件概览

`main.go` 是 Sentinel 项目的应用入口，负责初始化所有核心组件并启动 HTTP 服务器。

**文件位置**：`main.go`
**监听端口**：6872
**Web 框架**：GoFrame v2.7.1

---

## 完整代码结构

```go
package main

import (
    _ "SuperBizAgent/internal/ai/skills/builtin"  // 空导入，注册内置技能
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
        group.Bind(chat.NewV1())
        group.Bind(subscription.NewV1())
        group.Bind(report.New())
        group.Bind(report.NewTemplateController())
        group.Bind(event.New())
        group.Bind(skill.NewV1())
    })

    // 8. 设置端口并启动
    s.SetPort(6872)
    s.Run()
}
```

---

## 启动流程详解

### 1. 上下文初始化（第 21 行）

```go
ctx := gctx.New()
```

**作用**：
- 创建 GoFrame 的上下文对象
- 用于在整个应用生命周期中传递请求上下文、日志、链路追踪等信息
- 这是 GoFrame 的核心设计，所有操作都基于 context

**为什么需要 context**：
- 请求链路追踪
- 超时控制
- 取消信号传递
- 日志关联

---

### 2. 配置文件读取（第 22-26 行）

```go
fileDir, err := g.Cfg().Get(ctx, "file_dir")
if err != nil {
    panic(err)
}
common.FileDir = fileDir.String()
```

**配置管理器 `g.Cfg()`**：
- 单例模式，全局唯一
- 默认读取 `config/config.yaml` 或 `manifest/config/config.yaml`
- 支持多环境配置（dev、test、prod）
- 支持配置热更新

**配置文件搜索路径**：
1. `./config/config.yaml`
2. `./manifest/config/config.yaml`
3. 环境变量覆盖

**这里的作用**：
- 读取文件存储目录配置
- 保存到全局变量 `common.FileDir`
- 用于后续文件上传、下载等操作

---

### 3. 数据库初始化（第 28-31 行）

```go
if err := database.Init(); err != nil {
    panic(err)
}
```

**初始化内容**：
- 连接 MySQL 数据库
- 创建连接池（最大连接数、空闲连接数、超时时间）
- 执行数据库迁移（自动创建表结构）
- 连接 Milvus 向量数据库
- 创建 Milvus 集合（Collection）

**相关文件**：
- `internal/database/database.go` - 数据库初始化逻辑
- `internal/model/` - GORM 数据模型定义

---

### 4. 启动后台服务（第 33-34 行）

```go
service.InitFetcher()
```

**启动内容**：
- 初始化数据采集服务
- 启动定时任务调度器（cron）
- 加载所有订阅源配置
- 开始执行数据抓取任务

**相关文件**：
- `internal/service/fetcher.go` - 采集服务初始化
- `internal/scheduler/scheduler.go` - 定时任务调度
- `internal/fetcher/` - 各种数据源采集器

---

### 5. 获取 HTTP 服务器实例（第 36 行）

```go
s := g.Server()
```

**这是 GoFrame 的核心**：
- `g.Server()` 返回一个单例的 HTTP 服务器实例
- 基于 `net/http` 封装
- 支持优雅关闭、热重启
- 自动处理静态文件、模板渲染、Session 等

**单例模式的好处**：
- 全局唯一，避免重复创建
- 统一管理路由、中间件、配置
- 支持多端口监听（可以创建多个命名服务器）

---

### 6. 注册 SSE 路由（第 37-38 行）

```go
s.BindHandler("POST:/api/event/pipeline/stream", event.New().PipelineStream)
```

**为什么单独注册**：
- SSE（Server-Sent Events）需要保持长连接
- 不能经过响应格式化中间件（会破坏流式输出）
- 不需要 CORS 中间件（前端通过 EventSource 连接）

**SSE 的作用**：
- 实时推送 5 个 Agent 的执行过程
- 流式输出 AI 思考链
- 前端通过 `EventSource` 接收数据

**路由格式**：
- `POST:/api/event/pipeline/stream` - HTTP 方法 + 路径
- 直接绑定到控制器方法

---

### 7. 路由分组与中间件（第 40-49 行）

```go
s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(middleware.CORSMiddleware)
    group.Middleware(middleware.ResponseMiddleware)
    group.Bind(chat.NewV1())
    group.Bind(subscription.NewV1())
    group.Bind(report.New())
    group.Bind(report.NewTemplateController())
    group.Bind(event.New())
    group.Bind(skill.NewV1())
})
```

#### 7.1 路由分组

**`s.Group("/api", ...)`**：
- 创建路由组，所有路由都有 `/api` 前缀
- 支持嵌套分组（如 `/api/v1`、`/api/v2`）
- 可以为每个分组设置独立的中间件

#### 7.2 中间件注册

**执行顺序**：
1. `CORSMiddleware` - 处理跨域请求（OPTIONS 预检、响应头设置）
2. `ResponseMiddleware` - 统一响应格式（成功/失败）

**中间件链**：
```
请求 → CORSMiddleware → ResponseMiddleware → 控制器方法 → 响应
```

#### 7.3 控制器绑定

**`group.Bind()` 自动路由**：
- 自动扫描控制器的方法
- 根据方法名生成 RESTful 路由
- 支持参数验证、响应格式化

**示例**：
```go
type ChatController struct{}

// 自动生成路由：POST /api/chat/v1/create
func (c *ChatController) Create(ctx context.Context, req *v1.ChatCreateReq) (res *v1.ChatCreateRes, err error) {
    // 业务逻辑
}

// 自动生成路由：GET /api/chat/v1/list
func (c *ChatController) List(ctx context.Context, req *v1.ChatListReq) (res *v1.ChatListRes, err error) {
    // 业务逻辑
}
```

**绑定的控制器**：
- `chat.NewV1()` - 聊天接口（`/api/chat/v1/*`）
- `subscription.NewV1()` - 订阅管理（`/api/subscription/v1/*`）
- `report.New()` - 报告管理（`/api/report/*`）
- `report.NewTemplateController()` - 报告模板（`/api/report/template/*`）
- `event.New()` - 事件管理（`/api/event/*`）
- `skill.NewV1()` - 技能执行（`/api/skill/v1/*`）

---

### 8. 设置端口并启动（第 50-51 行）

```go
s.SetPort(6872)
s.Run()
```

**`s.SetPort(6872)`**：
- 设置监听端口为 6872
- 可以通过配置文件覆盖
- 支持多端口监听

**`s.Run()`**：
- 阻塞运行，直到收到退出信号
- 支持优雅关闭（Graceful Shutdown）
- 监听 `SIGTERM`、`SIGINT` 信号

**优雅关闭流程**：
1. 收到退出信号
2. 停止接收新请求
3. 等待现有请求处理完成（默认超时 30 秒）
4. 关闭数据库连接
5. 退出进程

---

## GoFrame 核心特性

### 1. 单例模式

```go
g.Server()  // 全局唯一的 HTTP 服务器
g.Cfg()     // 全局唯一的配置管理器
g.DB()      // 全局唯一的数据库连接池
g.Redis()   // 全局唯一的 Redis 客户端
```

**优点**：
- 避免重复创建
- 统一管理
- 线程安全

### 2. 自动路由绑定

**规范化路由**：
- 控制器方法名自动转换为路由路径
- 支持 RESTful 风格
- 自动参数验证

**示例**：
```go
// 方法名：Create → 路由：POST /create
// 方法名：GetList → 路由：GET /get-list
// 方法名：UpdateById → 路由：PUT /update-by-id
```

### 3. 配置文件自动加载

**搜索路径**：
1. `./config/config.yaml`
2. `./manifest/config/config.yaml`
3. 环境变量覆盖

**多环境支持**：
- `config.dev.yaml` - 开发环境
- `config.test.yaml` - 测试环境
- `config.prod.yaml` - 生产环境

### 4. 优雅关闭

**关闭流程**：
```
收到信号 → 停止接收新请求 → 等待现有请求完成 → 关闭连接 → 退出
```

**超时控制**：
- 默认等待 30 秒
- 可通过配置文件调整
- 超时后强制退出

---

## 启动流程图

```
main() 启动
    ↓
创建上下文 (gctx.New)
    ↓
读取配置文件 (g.Cfg)
    ↓
初始化数据库 (MySQL + Milvus)
    ↓
启动后台服务 (定时任务调度器)
    ↓
获取 HTTP 服务器 (g.Server)
    ↓
注册 SSE 路由 (不经过中间件)
    ↓
注册路由组 (/api)
    ├─ 注册中间件 (CORS + Response)
    └─ 绑定控制器 (自动扫描方法)
    ↓
设置端口 (6872)
    ↓
启动服务器 (s.Run) - 阻塞运行
    ↓
等待退出信号 (优雅关闭)
```

---

## 关键点总结

| 行号 | 代码 | 说明 |
|------|------|------|
| 4 | `_ "SuperBizAgent/internal/ai/skills/builtin"` | 空导入，执行包的 `init()` 函数，注册内置技能 |
| 21 | `ctx := gctx.New()` | 创建上下文，用于整个应用生命周期 |
| 22 | `g.Cfg().Get(ctx, "file_dir")` | 读取配置文件，支持多环境 |
| 29 | `database.Init()` | 初始化 MySQL + Milvus |
| 34 | `service.InitFetcher()` | 启动数据采集服务 |
| 36 | `s := g.Server()` | 获取单例 HTTP 服务器 |
| 38 | `s.BindHandler(...)` | 注册 SSE 路由，不经过中间件 |
| 40-49 | `s.Group("/api", ...)` | 路由分组 + 中间件 + 控制器绑定 |
| 50 | `s.SetPort(6872)` | 设置监听端口 |
| 51 | `s.Run()` | 阻塞运行，支持优雅关闭 |

---

## 相关文件

- `manifest/config/config.yaml` - 配置文件
- `internal/database/database.go` - 数据库初始化
- `internal/service/fetcher.go` - 数据采集服务
- `internal/controller/` - 控制器目录
- `utility/middleware/` - 中间件目录

---

## 下一步学习

1. 阅读 `internal/database/database.go` - 理解数据库初始化
2. 阅读 `internal/controller/event/event.go` - 理解控制器实现
3. 阅读 `utility/middleware/` - 理解中间件机制
4. 阅读 `manifest/config/config.yaml` - 理解配置结构