# 异常报警与自动召回功能使用说明

## 概述

在smartcontract.go中新增了异常报警与自动召回功能，用于检测农产品溯源过程中的异常情况并自动或手动触发产品召回。

## 新增功能

### 1. 异常检测 (DetectAnomalies)
- 自动检测过期产品
- 检测时间逻辑异常（生产时间早于采摘时间）
- 检测运输信息异常（司机年龄等）
- 检测联系方式异常

### 2. 异常报警 (AlertRecord)
- 创建异常报警记录
- 支持不同异常等级（1-5级，5级最严重）
- 跟踪报警处理状态

### 3. 自动召回 (AutoRecall)
- 当异常等级达到4级或以上时自动触发
- 生成召回记录并跟踪处理状态

### 4. 手动召回 (CreateRecall)
- 支持手动创建召回记录
- 可指定召回原因、等级和范围

### 5. 查询功能
- `GetFruitAlerts` - 获取指定产品的所有报警记录
- `GetFruitRecalls` - 获取指定产品的所有召回记录
- `GetPendingAlerts` - 获取所有待处理的报警
- `GetActiveRecalls` - 获取所有进行中的召回

## 使用方法

### 1. 上链时自动检测
```go
// 使用新的上链函数，会自动进行异常检测
txid, err := contract.UplinkWithDetection(ctx, userID, traceabilityCode, arg1, arg2, arg3, arg4, arg5, arg6)
```

### 2. 手动异常检测
```go
// 对指定产品进行异常检测
err := contract.DetectAnomalies(ctx, traceabilityCode)
```

### 3. 创建手动报警
```go
// 创建异常报警
err := contract.CreateAlert(ctx, traceabilityCode, "质量问题", "产品存在质量问题", 3, "operator123")
```

### 4. 创建手动召回
```go
// 创建手动召回
err := contract.CreateRecall(ctx, traceabilityCode, "质量问题召回", 4, "召回某批次产品", "operator123")
```

### 5. 更新报警/召回状态
```go
// 更新报警状态
err := contract.UpdateAlertStatus(ctx, alertID, "resolved", "问题已解决", "operator123")

// 更新召回状态
err := contract.UpdateRecallStatus(ctx, recallID, "completed", "召回已完成", "operator123")
```

### 6. 查询报警和召回记录
```go
// 获取产品的报警记录
alerts, err := contract.GetFruitAlerts(ctx, traceabilityCode)

// 获取产品的召回记录
recalls, err := contract.GetFruitRecalls(ctx, traceabilityCode)

// 获取所有待处理报警
pendingAlerts, err := contract.GetPendingAlerts(ctx)

// 获取所有进行中的召回
activeRecalls, err := contract.GetActiveRecalls(ctx)
```

## 异常等级说明

- **1级**: 轻微异常，记录即可
- **2级**: 一般异常，需要关注
- **3级**: 重要异常，需要处理
- **4级**: 严重异常，自动触发召回
- **5级**: 特别严重，立即处理

## 触发条件

### 自动触发召回的条件：
1. 生产时间逻辑异常（生产早于采摘）
2. 其他4级或5级异常

### 自动检测的异常类型：
1. **过期检测**: 销售时间已过期
2. **时间逻辑检测**: 生产时间早于采摘时间
3. **运输信息检测**: 司机年龄异常（<18岁或>70岁）
4. **联系信息检测**: 工厂联系方式异常

## 数据存储

- 报警记录存储为 `ALERT_<alertID>`
- 召回记录存储为 `RECALL_<recallID>`
- 支持通过范围查询获取所有记录

## 安全注意

- 异常检测和自动召回功能会在后台自动执行
- 关键操作（如召回）会有详细的日志记录
- 支持操作员信息追踪，便于审计