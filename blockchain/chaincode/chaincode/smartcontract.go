package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 定义合约结构体
type SmartContract struct {
	contractapi.Contract
}

// 注册用户
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, userID string, userType string, realInfoHash string) error {
	user := User{
		UserID:       userID,
		UserType:     userType,
		RealInfoHash: realInfoHash,
		FruitList:    []*Fruit{},
	}
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(userID, userAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// 农产品上链，传入用户ID、农产品上链信息
func (s *SmartContract) Uplink(ctx contractapi.TransactionContextInterface, userID string, traceability_code string, arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string) (string, error) {
	// 获取用户类型
	userType, err := s.GetUserType(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user type: %v", err)
	}

	// 通过溯源码获取农产品的上链信息
	FruitAsBytes, err := ctx.GetStub().GetState(traceability_code)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	// 将农产品的信息转换为Fruit结构体
	var fruit Fruit
	if FruitAsBytes != nil {
		err = json.Unmarshal(FruitAsBytes, &fruit)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal fruit: %v", err)
		}
	}

	//获取时间戳,修正时区
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to read TxTimestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai") // 选择你所在的时区
	time := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	// 获取交易ID
	txid := ctx.GetStub().GetTxID()
	// 给农产品信息中加上溯源码
	fruit.Traceability_code = traceability_code
	// 不同用户类型的上链的参数不一致
	switch userType {
	// 种植户
	case "种植户":
		// 将传入的农产品上链信息转换为Farmer_input结构体
		fruit.Farmer_input.Fa_fruitName = arg1
		fruit.Farmer_input.Fa_origin = arg2
		fruit.Farmer_input.Fa_plantTime = arg3
		fruit.Farmer_input.Fa_pickingTime = arg4
		fruit.Farmer_input.Fa_farmerName = arg5
		fruit.Farmer_input.Fa_img = arg6
		fruit.Farmer_input.Fa_Txid = txid
		fruit.Farmer_input.Fa_Timestamp = time
	// 工厂
	case "工厂":
		// 将传入的农产品上链信息转换为Factory_input结构体
		fruit.Factory_input.Fac_productName = arg1
		fruit.Factory_input.Fac_productionbatch = arg2
		fruit.Factory_input.Fac_productionTime = arg3
		fruit.Factory_input.Fac_factoryName = arg4
		fruit.Factory_input.Fac_contactNumber = arg5
		fruit.Factory_input.Fac_img = arg6
		fruit.Factory_input.Fac_Txid = txid
		fruit.Factory_input.Fac_Timestamp = time
	// 运输司机
	case "运输司机":
		// 将传入的农产品上链信息转换为Driver_input结构体
		fruit.Driver_input.Dr_name = arg1
		fruit.Driver_input.Dr_age = arg2
		fruit.Driver_input.Dr_phone = arg3
		fruit.Driver_input.Dr_carNumber = arg4
		fruit.Driver_input.Dr_transport = arg5
		fruit.Driver_input.Dr_img = arg6
		fruit.Driver_input.Dr_Txid = txid
		fruit.Driver_input.Dr_Timestamp = time
	// 商店
	case "商店":
		// 将传入的农产品上链信息转换为Shop_input结构体
		fruit.Shop_input.Sh_storeTime = arg1
		fruit.Shop_input.Sh_sellTime = arg2
		fruit.Shop_input.Sh_shopName = arg3
		fruit.Shop_input.Sh_shopAddress = arg4
		fruit.Shop_input.Sh_shopPhone = arg5
		fruit.Shop_input.Sh_img = arg6
		fruit.Shop_input.Sh_Txid = txid
		fruit.Shop_input.Sh_Timestamp = time

	}

	//将农产品的信息转换为json格式
	fruitAsBytes, err := json.Marshal(fruit)
	if err != nil {
		return "", fmt.Errorf("failed to marshal fruit: %v", err)
	}
	//将农产品的信息存入区块链
	err = ctx.GetStub().PutState(traceability_code, fruitAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to put fruit: %v", err)
	}
	//将农产品存入用户的农产品列表
	err = s.AddFruit(ctx, userID, &fruit)
	if err != nil {
		return "", fmt.Errorf("failed to add fruit to user: %v", err)

	}

	return txid, nil
}

// 农产品上链，传入用户ID、农产品上链信息
func (s *SmartContract) UplinkWithDetection(ctx contractapi.TransactionContextInterface, userID string, traceability_code string, arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string) (string, error) {
	// 先执行正常的上链操作
	txid, err := s.Uplink(ctx, userID, traceability_code, arg1, arg2, arg3, arg4, arg5, arg6)
	if err != nil {
		return "", err
	}

	// 执行异常检测
	err = s.DetectAnomalies(ctx, traceability_code)
	if err != nil {
		log.Printf("Anomaly detection failed for %s: %v", traceability_code, err)
		// 不返回错误，因为上链已经成功
	}

	return txid, nil
}

// 添加农产品到用户的农产品列表
func (s *SmartContract) AddFruit(ctx contractapi.TransactionContextInterface, userID string, fruit *Fruit) error {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return err
	}
	// 遍历用户的农产品列表，检查是否已经存在该农产品
	for _, existingFruit := range user.FruitList {
		if existingFruit.Traceability_code == fruit.Traceability_code {
			return fmt.Errorf("the fruit with traceability code %s already exists in user %s's fruit list", fruit.Traceability_code, userID)
		}
	}
	// 如果不存在，则将农产品添加到用户的农产品列表中
	user.FruitList = append(user.FruitList, fruit)
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(userID, userAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// 获取用户类型
func (s *SmartContract) GetUserType(ctx contractapi.TransactionContextInterface, userID string) (string, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return "", fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return "", err
	}
	return user.UserType, nil
}

// 获取用户信息
func (s *SmartContract) GetUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return &User{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return &User{}, fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

// 获取农产品的上链信息
func (s *SmartContract) GetFruitInfo(ctx contractapi.TransactionContextInterface, traceability_code string) (*Fruit, error) {
	FruitAsBytes, err := ctx.GetStub().GetState(traceability_code)
	if err != nil {
		return &Fruit{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// 将返回的结果转换为Fruit结构体
	var fruit Fruit
	if FruitAsBytes != nil {
		err = json.Unmarshal(FruitAsBytes, &fruit)
		if err != nil {
			return &Fruit{}, fmt.Errorf("failed to unmarshal fruit: %v", err)
		}
		if fruit.Traceability_code != "" {
			return &fruit, nil
		}
	}
	return &Fruit{}, fmt.Errorf("the fruit %s does not exist", traceability_code)
}

// 获取用户的农产品ID列表
func (s *SmartContract) GetFruitList(ctx contractapi.TransactionContextInterface, userID string) ([]*Fruit, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return user.FruitList, nil
}

// 获取所有的农产品信息
func (s *SmartContract) GetAllFruitInfo(ctx contractapi.TransactionContextInterface) ([]Fruit, error) {
	fruitListAsBytes, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	defer fruitListAsBytes.Close()
	var fruits []Fruit
	for fruitListAsBytes.HasNext() {
		queryResponse, err := fruitListAsBytes.Next()
		if err != nil {
			return nil, err
		}
		var fruit Fruit
		err = json.Unmarshal(queryResponse.Value, &fruit)
		if err != nil {
			return nil, err
		}
		// 过滤非农产品的信息
		if fruit.Traceability_code != "" {
			fruits = append(fruits, fruit)
		}
	}
	return fruits, nil
}

// 获取农产品上链历史
func (s *SmartContract) GetFruitHistory(ctx contractapi.TransactionContextInterface, traceability_code string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", traceability_code)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(traceability_code)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var fruit Fruit
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &fruit)
			if err != nil {
				return nil, err
			}
		} else {
			fruit = Fruit{
				Traceability_code: traceability_code,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}
		// 指定目标时区
		targetLocation, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return nil, err
		}

		// 将时间戳转换到目标时区
		timestamp = timestamp.In(targetLocation)
		// 格式化时间戳为指定格式的字符串
		formattedTime := timestamp.Format("2006-01-02 15:04:05")

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &fruit,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// 异常检测函数
func (s *SmartContract) DetectAnomalies(ctx contractapi.TransactionContextInterface, traceability_code string) error {
	fruit, err := s.GetFruitInfo(ctx, traceability_code)
	if err != nil {
		return fmt.Errorf("failed to get fruit info: %v", err)
	}

	// 检测各种异常情况
	anomalies := []string{}
	alertLevel := 0

	// 1. 检测过期时间异常
	if fruit.Shop_input.Sh_sellTime != "" {
		sellTime, err := time.Parse("2006-01-02 15:04:05", fruit.Shop_input.Sh_sellTime)
		if err == nil && sellTime.Before(time.Now()) {
			anomalies = append(anomalies, "农产品已过期")
			alertLevel = 3
		}
	}

	// 2. 检测生产时间逻辑异常
	if fruit.Farmer_input.Fa_pickingTime != "" && fruit.Factory_input.Fac_productionTime != "" {
		pickTime, err1 := time.Parse("2006-01-02 15:04:05", fruit.Farmer_input.Fa_pickingTime)
		prodTime, err2 := time.Parse("2006-01-02 15:04:05", fruit.Factory_input.Fac_productionTime)
		if err1 == nil && err2 == nil && prodTime.Before(pickTime) {
			anomalies = append(anomalies, "生产时间早于采摘时间，存在逻辑异常")
			if alertLevel < 4 {
				alertLevel = 4
			}
		}
	}

	// 3. 检测运输信息异常
	if fruit.Driver_input.Dr_age != "" {
		age, err := strconv.Atoi(fruit.Driver_input.Dr_age)
		if err == nil && (age < 18 || age > 70) {
			anomalies = append(anomalies, "司机年龄异常")
			if alertLevel < 2 {
				alertLevel = 2
			}
		}
	}

	// 4. 检测联系方式异常
	if fruit.Factory_input.Fac_contactNumber != "" && len(fruit.Factory_input.Fac_contactNumber) < 7 {
		anomalies = append(anomalies, "工厂联系方式异常")
		if alertLevel < 2 {
			alertLevel = 2
		}
	}

	// 如果有异常，创建报警记录
	if len(anomalies) > 0 {
		alertDesc := strings.Join(anomalies, "; ")
		err = s.CreateAlert(ctx, traceability_code, "自动检测", alertDesc, alertLevel, "system")
		if err != nil {
			return fmt.Errorf("failed to create alert: %v", err)
		}

		// 如果异常等级达到4级或以上，自动触发召回
		if alertLevel >= 4 {
			err = s.AutoRecall(ctx, traceability_code, alertDesc, alertLevel, "system")
			if err != nil {
				return fmt.Errorf("failed to auto recall: %v", err)
			}
		}
	}

	return nil
}

// 创建异常报警
func (s *SmartContract) CreateAlert(ctx contractapi.TransactionContextInterface, traceability_code string, alertType string, alertDesc string, alertLevel int, operatorID string) error {
	// 生成报警ID
	hash := sha256.Sum256([]byte(traceability_code + alertType + time.Now().String()))
	alertID := hex.EncodeToString(hash[:])[:16]

	// 获取当前时间
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	// 获取交易ID
	txid := ctx.GetStub().GetTxID()

	alert := AlertRecord{
		AlertID:          alertID,
		TraceabilityCode: traceability_code,
		AlertType:        alertType,
		AlertDesc:        alertDesc,
		AlertLevel:       alertLevel,
		AlertTime:        currentTime,
		Status:           "pending",
		Result:           "",
		OperatorID:       operatorID,
		Txid:             txid,
		Timestamp:        currentTime,
	}

	alertAsBytes, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %v", err)
	}

	// 存储报警记录
	err = ctx.GetStub().PutState("ALERT_"+alertID, alertAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put alert: %v", err)
	}

	log.Printf("Alert created: %s for fruit %s with level %d", alertID, traceability_code, alertLevel)
	return nil
}

// 自动召回功能
func (s *SmartContract) AutoRecall(ctx contractapi.TransactionContextInterface, traceability_code string, recallReason string, recallLevel int, operatorID string) error {
	// 生成召回ID
	hash := sha256.Sum256([]byte(traceability_code + "AUTO_RECALL" + time.Now().String()))
	recallID := hex.EncodeToString(hash[:])[:16]

	// 获取当前时间
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	// 获取交易ID
	txid := ctx.GetStub().GetTxID()

	recall := RecallRecord{
		RecallID:         recallID,
		TraceabilityCode: traceability_code,
		RecallReason:     recallReason,
		RecallLevel:      recallLevel,
		RecallTime:       currentTime,
		Status:           "initiated",
		ScopeDesc:        "自动召回 - 涉及该溯源码的所有相关产品",
		Result:           "自动召回已启动",
		OperatorID:       operatorID,
		Txid:             txid,
		Timestamp:        currentTime,
	}

	recallAsBytes, err := json.Marshal(recall)
	if err != nil {
		return fmt.Errorf("failed to marshal recall: %v", err)
	}

	// 存储召回记录
	err = ctx.GetStub().PutState("RECALL_"+recallID, recallAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put recall: %v", err)
	}

	log.Printf("Auto recall initiated: %s for fruit %s with level %d", recallID, traceability_code, recallLevel)
	return nil
}

// 手动创建召回记录
func (s *SmartContract) CreateRecall(ctx contractapi.TransactionContextInterface, traceability_code string, recallReason string, recallLevel int, scopeDesc string, operatorID string) error {
	// 生成召回ID
	hash := sha256.Sum256([]byte(traceability_code + "MANUAL_RECALL" + time.Now().String()))
	recallID := hex.EncodeToString(hash[:])[:16]

	// 获取当前时间
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	// 获取交易ID
	txid := ctx.GetStub().GetTxID()

	recall := RecallRecord{
		RecallID:         recallID,
		TraceabilityCode: traceability_code,
		RecallReason:     recallReason,
		RecallLevel:      recallLevel,
		RecallTime:       currentTime,
		Status:           "initiated",
		ScopeDesc:        scopeDesc,
		Result:           "手动召回已启动",
		OperatorID:       operatorID,
		Txid:             txid,
		Timestamp:        currentTime,
	}

	recallAsBytes, err := json.Marshal(recall)
	if err != nil {
		return fmt.Errorf("failed to marshal recall: %v", err)
	}

	// 存储召回记录
	err = ctx.GetStub().PutState("RECALL_"+recallID, recallAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put recall: %v", err)
	}

	log.Printf("Manual recall created: %s for fruit %s by operator %s", recallID, traceability_code, operatorID)
	return nil
}

// 更新报警状态
func (s *SmartContract) UpdateAlertStatus(ctx contractapi.TransactionContextInterface, alertID string, status string, result string, operatorID string) error {
	alertAsBytes, err := ctx.GetStub().GetState("ALERT_" + alertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %v", err)
	}
	if alertAsBytes == nil {
		return fmt.Errorf("alert %s does not exist", alertID)
	}

	var alert AlertRecord
	err = json.Unmarshal(alertAsBytes, &alert)
	if err != nil {
		return fmt.Errorf("failed to unmarshal alert: %v", err)
	}

	// 获取当前时间
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	alert.Status = status
	alert.Result = result
	alert.OperatorID = operatorID
	alert.Timestamp = currentTime

	alertAsBytes, err = json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %v", err)
	}

	err = ctx.GetStub().PutState("ALERT_"+alertID, alertAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put alert: %v", err)
	}

	log.Printf("Alert %s status updated to %s by operator %s", alertID, status, operatorID)
	return nil
}

// 更新召回状态
func (s *SmartContract) UpdateRecallStatus(ctx contractapi.TransactionContextInterface, recallID string, status string, result string, operatorID string) error {
	recallAsBytes, err := ctx.GetStub().GetState("RECALL_" + recallID)
	if err != nil {
		return fmt.Errorf("failed to get recall: %v", err)
	}
	if recallAsBytes == nil {
		return fmt.Errorf("recall %s does not exist", recallID)
	}

	var recall RecallRecord
	err = json.Unmarshal(recallAsBytes, &recall)
	if err != nil {
		return fmt.Errorf("failed to unmarshal recall: %v", err)
	}

	// 获取当前时间
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	recall.Status = status
	recall.Result = result
	recall.OperatorID = operatorID
	recall.Timestamp = currentTime

	recallAsBytes, err = json.Marshal(recall)
	if err != nil {
		return fmt.Errorf("failed to marshal recall: %v", err)
	}

	err = ctx.GetStub().PutState("RECALL_"+recallID, recallAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put recall: %v", err)
	}

	log.Printf("Recall %s status updated to %s by operator %s", recallID, status, operatorID)
	return nil
}

// 获取农产品的所有报警记录
func (s *SmartContract) GetFruitAlerts(ctx contractapi.TransactionContextInterface, traceability_code string) ([]AlertRecord, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("ALERT_", "ALERT_\xFF")
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %v", err)
	}
	defer resultsIterator.Close()

	var alerts []AlertRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var alert AlertRecord
		err = json.Unmarshal(queryResponse.Value, &alert)
		if err != nil {
			continue
		}

		if alert.TraceabilityCode == traceability_code {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// 获取农产品的所有召回记录
func (s *SmartContract) GetFruitRecalls(ctx contractapi.TransactionContextInterface, traceability_code string) ([]RecallRecord, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("RECALL_", "RECALL_\xFF")
	if err != nil {
		return nil, fmt.Errorf("failed to get recalls: %v", err)
	}
	defer resultsIterator.Close()

	var recalls []RecallRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var recall RecallRecord
		err = json.Unmarshal(queryResponse.Value, &recall)
		if err != nil {
			continue
		}

		if recall.TraceabilityCode == traceability_code {
			recalls = append(recalls, recall)
		}
	}

	return recalls, nil
}

// 获取所有待处理的报警
func (s *SmartContract) GetPendingAlerts(ctx contractapi.TransactionContextInterface) ([]AlertRecord, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("ALERT_", "ALERT_\xFF")
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %v", err)
	}
	defer resultsIterator.Close()

	var alerts []AlertRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var alert AlertRecord
		err = json.Unmarshal(queryResponse.Value, &alert)
		if err != nil {
			continue
		}

		if alert.Status == "pending" {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// 获取所有进行中的召回
func (s *SmartContract) GetActiveRecalls(ctx contractapi.TransactionContextInterface) ([]RecallRecord, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("RECALL_", "RECALL_\xFF")
	if err != nil {
		return nil, fmt.Errorf("failed to get recalls: %v", err)
	}
	defer resultsIterator.Close()

	var recalls []RecallRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var recall RecallRecord
		err = json.Unmarshal(queryResponse.Value, &recall)
		if err != nil {
			continue
		}

		if recall.Status == "initiated" || recall.Status == "processing" {
			recalls = append(recalls, recall)
		}
	}

	return recalls, nil
}
