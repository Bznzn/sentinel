# Sentinel 后端学习路径

## 🎯 学习目标

从0到精通 Sentinel 项目的后端系统，掌握 Go 后端开发、AI Agent 协作、RAG 技术和生产级部署。

预计学习时间：**3-4 周**

---

## 第一阶段：Go 后端基础架构（2-3天）

### 目标
理解项目的分层架构和启动流程

### 学习内容

#### 1. 应用启动流程
- 阅读 `main.go:37` - 应用入口
- 理解 GoFrame 的启动机制
- 查看路由注册：`internal/controller/`

#### 2. 分层架构理解
```
Controller（HTTP层）→ Service（业务层）→ Repository（数据层）→ Model（数据模型）
```

#### 3. 核心组件
- 配置管理：`internal/config/`
- 数据库初始化：`internal/database/`
- 中间件：`internal/middleware/`（CORS、Auth、RateLimit）
- 缓存：`internal/cache/`

### 关键文件
- `main.go` - 启动入口
- `manifest/config/config.yaml` - 配置文件
- `internal/model/` - GORM 数据模型
- `internal/repository/` - 数据访问层

### 实践任务
- [ ] 启动项目，访问 API 文档
- [ ] 修改配置文件，观察效果
- [ ] 添加一个简单的 API 端点（如 `/api/health`）

---

## 第二阶段：数据采集系统（3-4天）

### 目标
掌握多数据源采集和定时任务调度

### 学习内容

#### 1. 数据源采集器
```
internal/fetcher/
├── github.go      # GitHub 仓库采集
├── nvd.go         # NVD 漏洞数据库
├── rss.go         # RSS 订阅源
└── webhook.go     # Webhook 接收
```

#### 2. 定时任务调度
- 查看 `internal/scheduler/scheduler.go`
- 理解 cron 表达式配置
- 学习任务并发控制

#### 3. 外部 API 客户端
- GitHub API：`internal/client/github/`
- HTTP 客户端封装
- 错误处理和重试机制

### 核心流程
```
订阅管理 → 定时调度 → 数据采集 → 数据清洗 → 入库存储
```

### 实践任务
- [ ] 添加一个新的 RSS 订阅源
- [ ] 实现一个新的数据源采集器（如 Twitter API）
- [ ] 编写单元测试：`internal/fetcher/*_test.go`

---

## 第三阶段：AI Agent 核心系统（5-7天）

### 目标
深入理解多 Agent 协作和 RAG 技术

### 3.1 Cloudwego Eino 框架学习

先理解框架的核心概念：
- **Graph**：Agent 工作流图
- **Node**：单个处理节点
- **State**：状态传递
- **Stream**：流式输出

### 3.2 五个核心 Agent 解析

#### Agent 1: 智能提取
**位置**：`internal/ai/agent/event_extraction/`

- 使用 DeepSeek V3 结构化提取
- 输入：原始事件文本
- 输出：标题、CVE、严重程度、产品等

#### Agent 2: 语义去重
**位置**：`internal/ai/agent/event_dedup/`

- MD5 哈希去重
- 基于标题 + CVE 生成唯一标识
- 避免重复告警

#### Agent 3: 风险评估
**位置**：`internal/ai/agent/risk_assessment/`

- AI 评分（0-100）
- 输出 severity 和 recommendation
- 结合 CVSS 评分

#### Agent 4: 解决方案生成
**位置**：`internal/ai/agent/event_analysis_pipeline/`

- ReAct 推理循环（最多 25 步）
- 调用 6 个工具
- 生成处置方案

#### Agent 5: 向量索引
**位置**：`internal/ai/agent/event_indexer/`

- 事件向量化
- 写入 Milvus
- 容量管理（10000 条上限）

### 3.3 ReAct 工具系统

**位置**：`internal/ai/tools/` 和 `internal/ai/agent/chat_pipeline/tools_node.go`

6 个核心工具：
1. 查询历史事件
2. 搜索知识库
3. 查询 CVE 详情
4. 威胁情报查询
5. 日志分析
6. 网络搜索

### 3.4 RAG 知识库

#### 向量化
**位置**：`internal/ai/embedder/`

- 使用豆包 Embedding（1024 维）
- 文本分块和向量化

#### 索引管理
**位置**：`internal/ai/indexer/`

- Milvus 集合创建
- 数据写入和 Flush
- 容量监控

#### 检索
**位置**：`internal/ai/retriever/`

- 相似度搜索（TopK=3）
- 语义检索历史事件
- RAG 上下文增强

### 关键文件
- `internal/ai/agent/event_pipeline/event_pipeline.go` - 5 个 Agent 串联
- `internal/ai/agent/chat_pipeline/tools_node.go` - ReAct 工具调用
- `internal/ai/models/` - LLM 模型集成

### 实践任务
- [ ] 调试一次完整的 5 Agent 管道
- [ ] 添加一个新工具（如"查询威胁情报"）
- [ ] 修改风险评估的评分逻辑
- [ ] 优化向量检索的 TopK 参数

---

## 第四阶段：生产级优化（3-5天）

### 目标
性能优化、安全加固、部署运维

### 1. 性能优化
- 数据库查询优化（索引、分页）
- Milvus 向量检索调优
- 并发控制和连接池
- 缓存策略：`internal/cache/`

### 2. 安全加固
- JWT 认证：`internal/auth/`
- RBAC 权限控制
- API 限流：`internal/middleware/`
- 敏感信息加密

### 3. 监控与日志
- 健康检查：`internal/health/`
- 结构化日志
- 性能指标采集

### 4. Docker 部署
- 多阶段构建：`Dockerfile`
- Docker Compose 编排
- Milvus 集群部署：`manifest/docker/docker-compose.yml`

### 实践任务
- [ ] 添加数据库索引优化查询
- [ ] 实现 API 限流中间件
- [ ] 编写健康检查端点
- [ ] 优化 Docker 镜像大小

---

## 📖 推荐学习顺序

| 天数 | 学习内容 |
|------|---------|
| 第 1 天 | main.go → config.yaml → 数据模型 |
| 第 2-3 天 | Controller → Service → Repository 分层架构 |
| 第 4-5 天 | 数据采集器 → 定时调度 → 外部 API 客户端 |
| 第 6-7 天 | Cloudwego Eino 框架 → Agent 1-3（提取、去重、评估） |
| 第 8-10 天 | Agent 4-5（ReAct、索引）→ RAG 知识库 |
| 第 11-12 天 | 工具系统 → Milvus 向量数据库 |
| 第 13-15 天 | 性能优化 → 安全加固 → 部署运维 |

---

## 🔧 必备技能

1. **Go 语言**：基础语法、并发编程（goroutine、channel）
2. **GoFrame 框架**：路由、中间件、ORM
3. **GORM**：数据库操作、关联查询
4. **Milvus**：向量数据库基础
5. **LLM API**：DeepSeek、豆包 Embedding
6. **Docker**：容器化部署

---

## 📚 核心文档

- [GoFrame 官方文档](https://goframe.org/)
- [Cloudwego Eino GitHub](https://github.com/cloudwego/eino)
- [Milvus 官方文档](https://milvus.io/docs)
- [GORM 官方文档](https://gorm.io/)

---

## 💡 学习建议

1. **边学边做**：每个阶段都要动手实践，不要只看代码
2. **调试技巧**：善用断点调试，追踪完整的请求链路
3. **阅读源码**：重点关注 AI Agent 系统和 RAG 技术，这是项目的核心创新点
4. **提问交流**：遇到问题先查日志，再查文档，最后提问

---

## 🎯 学习重点

**核心创新点**：
- 多 Agent 协作架构（5 个 Agent 串联）
- ReAct 推理循环（Plan-Execute-Replan）
- RAG 知识库（Milvus 向量检索）
- SSE 流式输出（实时展示 Agent 思考过程）

专注于这些核心技术，理解它们如何协同工作，就能掌握整个项目的精髓。