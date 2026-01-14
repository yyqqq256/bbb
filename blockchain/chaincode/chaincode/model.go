package chaincode

/*
定义用户结构体
用户ID
用户类型
实名认证信息哈希,包括用户注册的姓名、身份证号、手机号、注册平台同意协议签名的哈希
农产品列表
*/
type User struct {
	UserID       string   `json:"userID"`
	UserType     string   `json:"userType"`
	RealInfoHash string   `json:"realInfoHash"`
	FruitList    []*Fruit `json:"fruitList"`
}

/*
定义农产品结构体
溯源码
种植户输入
工厂输入
运输司机输入
商店输入
*/
type Fruit struct {
	Traceability_code string        `json:"traceability_code"`
	Farmer_input      Farmer_input  `json:"farmer_input"`
	Factory_input     Factory_input `json:"factory_input"`
	Driver_input      Driver_input  `json:"driver_input"`
	Shop_input        Shop_input    `json:"shop_input"`
}

// HistoryQueryResult structure used for handling result of history query
type HistoryQueryResult struct {
	Record    *Fruit `json:"record"`
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
}

/*
种植户
农产品的溯源码，一物一码，主打高端市场（自动生成）
农产品名称
产地
种植时间
采摘时间
种植户名称
*/
type Farmer_input struct {
	Fa_fruitName   string `json:"fa_fruitName"`
	Fa_origin      string `json:"fa_origin"`
	Fa_plantTime   string `json:"fa_plantTime"`
	Fa_pickingTime string `json:"fa_pickingTime"`
	Fa_farmerName  string `json:"fa_farmerName"`
	Fa_img         string `json:"fa_img"` // 图片链接
	Fa_Txid        string `json:"fa_txid"`
	Fa_Timestamp   string `json:"fa_timestamp"`
}

/*
工厂
商品名称
生产批次
出厂时间（可以防止黑心商家修改时间）
工厂名称与厂址
联系电话
*/
type Factory_input struct {
	Fac_productName     string `json:"fac_productName"`
	Fac_productionbatch string `json:"fac_productionbatch"`
	Fac_productionTime  string `json:"fac_productionTime"`
	Fac_factoryName     string `json:"fac_factoryName"`
	Fac_contactNumber   string `json:"fac_contactNumber"`
	Fac_img             string `json:"fac_img"` // 图片链接
	Fac_Txid            string `json:"fac_txid"`
	Fac_Timestamp       string `json:"fac_timestamp"`
}

/*
运输司机
姓名
年龄
电话
车牌号
运输记录
*/
type Driver_input struct {
	Dr_name      string `json:"dr_name"`
	Dr_age       string `json:"dr_age"`
	Dr_phone     string `json:"dr_phone"`
	Dr_carNumber string `json:"dr_carNumber"`
	Dr_transport string `json:"dr_transport"`
	Dr_img       string `json:"dr_img"` // 图片链接
	Dr_Txid      string `json:"dr_txid"`
	Dr_Timestamp string `json:"dr_timestamp"`
}

/*
商店
存入时间
销售时间
商店名称
商店位置
商店电话
*/
type Shop_input struct {
	Sh_storeTime   string `json:"sh_storeTime"`
	Sh_sellTime    string `json:"sh_sellTime"`
	Sh_shopName    string `json:"sh_shopName"`
	Sh_shopAddress string `json:"sh_shopAddress"`
	Sh_shopPhone   string `json:"sh_shopPhone"`
	Sh_img         string `json:"sh_img"` // 图片链接
	Sh_Txid        string `json:"sh_txid"`
	Sh_Timestamp   string `json:"sh_timestamp"`
}

/*
异常报警记录
异常ID
溯源码
异常类型
异常描述
异常等级 (1-5, 5为最严重)
发现时间
处理状态
处理结果
操作员ID
*/
type AlertRecord struct {
	AlertID         string `json:"alert_id"`
	TraceabilityCode string `json:"traceability_code"`
	AlertType       string `json:"alert_type"`
	AlertDesc       string `json:"alert_desc"`
	AlertLevel      int    `json:"alert_level"`
	AlertTime       string `json:"alert_time"`
	Status          string `json:"status"` // pending, processing, resolved
	Result          string `json:"result"`
	OperatorID      string `json:"operator_id"`
	Txid            string `json:"txid"`
	Timestamp       string `json:"timestamp"`
}

/*
召回记录
召回ID
溯源码
召回原因
召回等级 (1-5, 5为最严重)
召回时间
召回状态
召回范围描述
处理结果
操作员ID
*/
type RecallRecord struct {
	RecallID         string `json:"recall_id"`
	TraceabilityCode string `json:"traceability_code"`
	RecallReason     string `json:"recall_reason"`
	RecallLevel      int    `json:"recall_level"`
	RecallTime       string `json:"recall_time"`
	Status           string `json:"status"` // initiated, processing, completed
	ScopeDesc        string `json:"scope_desc"`
	Result           string `json:"result"`
	OperatorID       string `json:"operator_id"`
	Txid             string `json:"txid"`
	Timestamp        string `json:"timestamp"`
}

/*
异常检测规则
规则ID
规则名称
检测类型
规则描述
触发条件
异常等级
是否启用
*/
type AlertRule struct {
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	DetectType  string `json:"detect_type"` // quality, safety, expiration, etc.
	RuleDesc    string `json:"rule_desc"`
	Condition   string `json:"condition"`
	AlertLevel  int    `json:"alert_level"`
	IsEnabled   bool   `json:"is_enabled"`
}
