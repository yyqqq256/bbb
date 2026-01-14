package controller

import (
	"backend/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 异常检测
func DetectAnomalies(c *gin.Context) {
	traceabilityCode := c.PostForm("traceability_code")
	if traceabilityCode == "" {
		c.JSON(200, gin.H{
			"message": "溯源码不能为空",
		})
		return
	}

	res, err := pkg.ChaincodeInvoke("DetectAnomalies", []string{traceabilityCode})
	if err != nil {
		c.JSON(200, gin.H{
			"message": "异常检测失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "异常检测完成",
		"txid":    res,
	})
}

// 创建异常报警
func CreateAlert(c *gin.Context) {
	traceabilityCode := c.PostForm("traceability_code")
	alertType := c.PostForm("alert_type")
	alertDesc := c.PostForm("alert_desc")
	alertLevelStr := c.PostForm("alert_level")
	operatorID := c.PostForm("operator_id")

	if traceabilityCode == "" || alertType == "" || alertDesc == "" || alertLevelStr == "" {
		c.JSON(200, gin.H{
			"message": "参数不完整",
		})
		return
	}

	alertLevel, err := strconv.Atoi(alertLevelStr)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "异常等级格式错误",
		})
		return
	}

	args := []string{traceabilityCode, alertType, alertDesc, alertLevelStr, operatorID}
	res, err := pkg.ChaincodeInvoke("CreateAlert", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "创建报警失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建报警成功",
		"txid":    res,
	})
}

// 创建召回记录
func CreateRecall(c *gin.Context) {
	traceabilityCode := c.PostForm("traceability_code")
	recallReason := c.PostForm("recall_reason")
	recallLevelStr := c.PostForm("recall_level")
	scopeDesc := c.PostForm("scope_desc")
	operatorID := c.PostForm("operator_id")

	if traceabilityCode == "" || recallReason == "" || recallLevelStr == "" {
		c.JSON(200, gin.H{
			"message": "参数不完整",
		})
		return
	}

	recallLevel, err := strconv.Atoi(recallLevelStr)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "召回等级格式错误",
		})
		return
	}

	if scopeDesc == "" {
		scopeDesc = "手动召回 - 涉及该溯源码的所有相关产品"
	}

	args := []string{traceabilityCode, recallReason, recallLevelStr, scopeDesc, operatorID}
	res, err := pkg.ChaincodeInvoke("CreateRecall", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "创建召回失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建召回成功",
		"txid":    res,
	})
}

// 更新报警状态
func UpdateAlertStatus(c *gin.Context) {
	alertID := c.PostForm("alert_id")
	status := c.PostForm("status")
	result := c.PostForm("result")
	operatorID := c.PostForm("operator_id")

	if alertID == "" || status == "" {
		c.JSON(200, gin.H{
			"message": "参数不完整",
		})
		return
	}

	args := []string{alertID, status, result, operatorID}
	res, err := pkg.ChaincodeInvoke("UpdateAlertStatus", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "更新报警状态失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新报警状态成功",
		"txid":    res,
	})
}

// 更新召回状态
func UpdateRecallStatus(c *gin.Context) {
	recallID := c.PostForm("recall_id")
	status := c.PostForm("status")
	result := c.PostForm("result")
	operatorID := c.PostForm("operator_id")

	if recallID == "" || status == "" {
		c.JSON(200, gin.H{
			"message": "参数不完整",
		})
		return
	}

	args := []string{recallID, status, result, operatorID}
	res, err := pkg.ChaincodeInvoke("UpdateRecallStatus", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "更新召回状态失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新召回状态成功",
		"txid":    res,
	})
}

// 获取产品的报警记录
func GetFruitAlerts(c *gin.Context) {
	traceabilityCode := c.PostForm("traceability_code")
	if traceabilityCode == "" {
		c.JSON(200, gin.H{
			"message": "溯源码不能为空",
		})
		return
	}

	res, err := pkg.ChaincodeQuery("GetFruitAlerts", traceabilityCode)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询报警记录失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "查询报警记录成功",
		"data":    res,
	})
}

// 获取产品的召回记录
func GetFruitRecalls(c *gin.Context) {
	traceabilityCode := c.PostForm("traceability_code")
	if traceabilityCode == "" {
		c.JSON(200, gin.H{
			"message": "溯源码不能为空",
		})
		return
	}

	res, err := pkg.ChaincodeQuery("GetFruitRecalls", traceabilityCode)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询召回记录失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "查询召回记录成功",
		"data":    res,
	})
}

// 获取所有待处理报警
func GetPendingAlerts(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetPendingAlerts", "")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询待处理报警失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "查询待处理报警成功",
		"data":    res,
	})
}

// 获取所有进行中的召回
func GetActiveRecalls(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetActiveRecalls", "")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询进行中召回失败：" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "查询进行中召回成功",
		"data":    res,
	})
}