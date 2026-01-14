package model

import (
	"time"
)

// 异常类型定义
type AlertType string

const (
	AlertType_Quality    AlertType = "quality"     // 质量问题
	AlertType_Expired    AlertType = "expired"     // 过期产品
	AlertType_Contaminated AlertType = "contaminated" // 污染产品
	AlertType_Fake       AlertType = "fake"        // 假冒产品
	AlertType_Recall     AlertType = "recall"      // 主动召回
	AlertType_Other      AlertType = "other"       // 其他异常
)

// 异常级别定义
type AlertLevel string

const (
	AlertLevel_Low      AlertLevel = "low"      // 低风险
	AlertLevel_Medium   AlertLevel = "medium"   // 中等风险
	AlertLevel_High     AlertLevel = "high"     // 高风险
	AlertLevel_Critical AlertLevel = "critical" // 严重风险
)

// 处理状态定义
type AlertStatus string

const (
	AlertStatus_Pending   AlertStatus = "pending"   // 待处理
	AlertStatus_Processing AlertStatus = "processing" // 处理中
	AlertStatus_Resolved AlertStatus = "resolved"  // 已解决
	AlertStatus_Closed   AlertStatus = "closed"    // 已关闭
)

// 召回状态定义
type RecallStatus string

const (
	RecallStatus_Draft    RecallStatus = "draft"     // 草稿
	RecallStatus_Active   RecallStatus = "active"    // 执行中
	RecallStatus_Complete RecallStatus = "complete"  // 已完成
	RecallStatus_Cancelled RecallStatus = "cancelled" // 已取消
)

// 异常报警记录
type Alert struct {
	ID                uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	AlertCode         string       `json:"alert_code" gorm:"uniqueIndex;not null"` // 异常编号
	TraceabilityCode  string       `json:"traceability_code" gorm:"index;not null"` // 溯源码
	ProductName       string       `json:"product_name" gorm:"not null"`          // 产品名称
	AlertType         AlertType    `json:"alert_type" gorm:"not null"`            // 异常类型
	AlertLevel        AlertLevel   `json:"alert_level" gorm:"not null"`           // 异常级别
	Title             string       `json:"title" gorm:"not null"`                   // 异常标题
	Description       string       `json:"description" gorm:"type:text"`           // 异常描述
	EvidenceImages    string       `json:"evidence_images" gorm:"type:text"`     // 证据图片(JSON数组)
	ReporterID        string       `json:"reporter_id" gorm:"index;not null"`     // 报告人ID
	ReporterName      string       `json:"reporter_name" gorm:"not null"`         // 报告人姓名
	ReporterContact   string       `json:"reporter_contact" gorm:"not null"`      // 报告人联系方式
	Status            AlertStatus  `json:"status" gorm:"default:pending"`          // 处理状态
	AssignedTo        string       `json:"assigned_to" gorm:"index"`              // 指派给
	HandlingMeasures  string       `json:"handling_measures" gorm:"type:text"`    // 处理措施
	HandlingResult    string       `json:"handling_result" gorm:"type:text"`      // 处理结果
	HandlingTime      *time.Time   `json:"handling_time"`                           // 处理时间
	RecallInitiated   bool         `json:"recall_initiated" gorm:"default:false"`  // 是否已发起召回
	RecallCode        string       `json:"recall_code" gorm:"index"`              // 关联的召回编号
	CreatedAt         time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// 产品召回记录
type ProductRecall struct {
	ID                uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	RecallCode        string        `json:"recall_code" gorm:"uniqueIndex;not null"` // 召回编号
	AlertCode         string        `json:"alert_code" gorm:"index"`                 // 关联的异常编号
	TraceabilityCode  string        `json:"traceability_code" gorm:"index;not null"` // 溯源码
	ProductName       string        `json:"product_name" gorm:"not null"`            // 产品名称
	BatchNumber       string        `json:"batch_number" gorm:"index"`                 // 批次号
	RecallReason      string        `json:"recall_reason" gorm:"type:text"`          // 召回原因
	RecallLevel       AlertLevel    `json:"recall_level" gorm:"not null"`             // 召回级别
	RecallScope       string        `json:"recall_scope" gorm:"type:text"`           // 召回范围
	RecallInstructions string       `json:"recall_instructions" gorm:"type:text"`   // 召回说明
	InitiatorID       string        `json:"initiator_id" gorm:"not null"`            // 发起人ID
	InitiatorName     string        `json:"initiator_name" gorm:"not null"`          // 发起人姓名
	InitiatorContact  string        `json:"initiator_contact" gorm:"not null"`       // 发起人联系方式
	Status            RecallStatus  `json:"status" gorm:"default:draft"`             // 召回状态
	ApprovedBy        string        `json:"approved_by" gorm:"index"`                  // 审批人
	ApprovalTime      *time.Time    `json:"approval_time"`                            // 审批时间
	ApprovalComment   string        `json:"approval_comment" gorm:"type:text"`        // 审批意见
	EffectiveTime     *time.Time    `json:"effective_time"`                             // 生效时间
	CompletionTime    *time.Time    `json:"completion_time"`                            // 完成时间
	AffectedQuantity  int           `json:"affected_quantity" gorm:"default:0"`       // 影响数量
	RecalledQuantity  int           `json:"recalled_quantity" gorm:"default:0"`       // 已召回数量
	RecallProgress    float64       `json:"recall_progress" gorm:"default:0"`         // 召回进度(百分比)
	PublicNotice      string        `json:"public_notice" gorm:"type:text"`           // 公告内容
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// 召回处理记录
type RecallProcess struct {
	ID                uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	RecallCode        string        `json:"recall_code" gorm:"index;not null"`       // 召回编号
	ProcessType       string        `json:"process_type" gorm:"not null"`            // 处理类型
	Description       string        `json:"description" gorm:"type:text"`              // 处理描述
	ProcessedBy       string        `json:"processed_by" gorm:"not null"`              // 处理人
	ProcessedByName   string        `json:"processed_by_name" gorm:"not null"`         // 处理人姓名
	ProcessedQuantity int           `json:"processed_quantity" gorm:"default:0"`       // 处理数量
	ProcessTime       time.Time     `json:"process_time" gorm:"autoCreateTime"`
	Location          string        `json:"location" gorm:"not null"`                  // 处理地点
	EvidenceImages    string        `json:"evidence_images" gorm:"type:text"`         // 证据图片
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
}

// 异常检测规则
type AlertRule struct {
	ID                uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	RuleCode          string        `json:"rule_code" gorm:"uniqueIndex;not null"`     // 规则编号
	RuleName          string        `json:"rule_name" gorm:"not null"`                   // 规则名称
	Description       string        `json:"description" gorm:"type:text"`                 // 规则描述
	AlertType         AlertType     `json:"alert_type" gorm:"not null"`                    // 异常类型
	AlertLevel        AlertLevel    `json:"alert_level" gorm:"not null"`                    // 异常级别
	TriggerConditions string        `json:"trigger_conditions" gorm:"type:text"`           // 触发条件(JSON)
	DetectionScope    string        `json:"detection_scope" gorm:"type:text"`             // 检测范围
	IsEnabled         bool          `json:"is_enabled" gorm:"default:true"`                 // 是否启用
	CreatedBy         string        `json:"created_by" gorm:"not null"`                     // 创建人
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// 检测异常请求
type DetectAnomalyRequest struct {
	TraceabilityCode string                 `json:"traceability_code" binding:"required"` // 溯源码
	ProductData      map[string]interface{} `json:"product_data" binding:"required"`     // 产品数据
	DetectionRules   []string               `json:"detection_rules"`                       // 指定检测规则
}

// 创建异常请求
type CreateAlertRequest struct {
	TraceabilityCode string    `json:"traceability_code" binding:"required"` // 溯源码
	AlertType        AlertType `json:"alert_type" binding:"required"`        // 异常类型
	AlertLevel       AlertLevel `json:"alert_level" binding:"required"`        // 异常级别
	Title            string    `json:"title" binding:"required"`             // 异常标题
	Description      string    `json:"description" binding:"required"`       // 异常描述
	EvidenceImages   []string  `json:"evidence_images"`                      // 证据图片
}

// 更新异常状态请求
type UpdateAlertStatusRequest struct {
	Status           AlertStatus `json:"status" binding:"required"`            // 状态
	HandlingMeasures string      `json:"handling_measures"`                    // 处理措施
	HandlingResult   string      `json:"handling_result"`                      // 处理结果
	AssignedTo       string      `json:"assigned_to"`                          // 指派给
}

// 创建召回请求
type CreateRecallRequest struct {
	TraceabilityCode   string     `json:"traceability_code" binding:"required"` // 溯源码
	AlertCode          string     `json:"alert_code"`                            // 关联的异常编号
	BatchNumber        string     `json:"batch_number"`                          // 批次号
	RecallReason       string     `json:"recall_reason" binding:"required"`      // 召回原因
	RecallLevel        AlertLevel `json:"recall_level" binding:"required"`       // 召回级别
	RecallScope        string     `json:"recall_scope" binding:"required"`       // 召回范围
	RecallInstructions string     `json:"recall_instructions" binding:"required"` // 召回说明
	PublicNotice       string     `json:"public_notice"`                          // 公告内容
}
