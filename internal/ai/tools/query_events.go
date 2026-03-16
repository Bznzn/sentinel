package tools

import (
	"context"
	"encoding/json"

	"SuperBizAgent/internal/database"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type EventQueryTool struct{}

func NewEventQueryTool() tool.BaseTool {
	return &EventQueryTool{}
}

func (t *EventQueryTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_events",
		Desc: `查询安全事件数据库。返回事件ID、标题、严重程度、状态和时间。

使用场景:
- 用户问"最近有什么安全事件" → 不传参数
- 用户问"有哪些高危漏洞" → severity="critical" 或 "high"
- 用户问"未处理的事件" → status="new"
- 用户问"CVE-2021-44228是什么漏洞" → cve_id="CVE-2021-44228"

返回格式: {"count": N, "items": [{id, title, severity, status, event_time, cve_id, description}]}`,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"severity": {Type: schema.String, Desc: "严重程度过滤。可选值: critical(紧急), high(高危), medium(中危), low(低危), info(信息)。不传则返回所有级别"},
			"status":   {Type: schema.String, Desc: "状态过滤。可选值: new(新建), processing(处理中), resolved(已解决), ignored(已忽略)。不传则返回所有状态"},
			"cve_id":   {Type: schema.String, Desc: "CVE编号，例如: CVE-2021-44228"},
			"limit":    {Type: schema.Integer, Desc: "返回数量，1-5之间，默认5"},
		}),
	}, nil
}

func (t *EventQueryTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params struct {
		Severity string `json:"severity"`
		Status   string `json:"status"`
		CveID    string `json:"cve_id"`
		Limit    int    `json:"limit"`
	}
	json.Unmarshal([]byte(args), &params)
	if params.Limit <= 0 || params.Limit > 5 {
		params.Limit = 5
	}

	db := database.GetDB()
	query := db.Table("security_events").Select("id, title, severity, status, event_time, cve_id, description, cvss_score, affected_vendor, affected_product").Order("event_time DESC").Limit(params.Limit)
	if params.Severity != "" {
		query = query.Where("severity = ?", params.Severity)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.CveID != "" {
		query = query.Where("cve_id = ?", params.CveID)
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return `{"error": "` + err.Error() + `"}`, nil
	}

	data, _ := json.Marshal(map[string]interface{}{"count": len(results), "items": results})
	return string(data), nil
}
