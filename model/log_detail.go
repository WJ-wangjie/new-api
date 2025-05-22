package model

type LogDetail struct {
	Id               int    `json:"id" gorm:"index:idx_created_at_id,priority:1"`
	UserId           int    `json:"user_id" gorm:"index"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;index:idx_created_at_id,priority:2;index:idx_created_at_type"`
	Type             int    `json:"type" gorm:"index:idx_created_at_type"`
	Content          string `json:"content"`
	Username         string `json:"username" gorm:"index;index:index_username_model_name,priority:2;default:''"`
	TokenName        string `json:"token_name" gorm:"index;default:''"`
	ModelName        string `json:"model_name" gorm:"index;index:index_username_model_name,priority:1;default:''"`
	Quota            int    `json:"quota" gorm:"default:0"`
	PromptTokens     int    `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int    `json:"completion_tokens" gorm:"default:0"`
	UseTime          int    `json:"use_time" gorm:"default:0"`
	IsStream         bool   `json:"is_stream" gorm:"default:false"`
	ChannelId        int    `json:"channel" gorm:"index"`
	ChannelName      string `json:"channel_name" gorm:"->"`
	TokenId          int    `json:"token_id" gorm:"default:0;index"`
	Group            string `json:"group" gorm:"index"`
	Other            string `json:"other"`
	RequestBody      string `json:"request_body"`
	ResponseBody     string `json:"response_body"`
}

//
//const (
//	LogDetailTypeUnknown = iota
//	LogDetailTypeTopup
//	LogDetailTypeConsume
//	LogDetailTypeManage
//	LogDetailTypeSystem
//	LogDetailTypeError
//)
//
//func formatUserLogDetails(logs []*LogDetail) {
//	for i := range logs {
//		logs[i].ChannelName = ""
//		var otherMap map[string]interface{}
//		otherMap = common.StrToMap(logs[i].Other)
//		if otherMap != nil {
//			// delete admin
//			delete(otherMap, "admin_info")
//		}
//		logs[i].Other = common.MapToJsonStr(otherMap)
//		logs[i].Id = logs[i].Id % 1024
//	}
//}
//
//func GetLogDetailByKey(key string) (logs []*LogDetail, err error) {
//	if os.Getenv("LOG_SQL_DSN") != "" {
//		var tk Token
//		if err = DB.Model(&Token{}).Where(keyCol+"=?", strings.TrimPrefix(key, "sk-")).First(&tk).Error; err != nil {
//			return nil, err
//		}
//		err = LOG_DB.Model(&LogDetail{}).Where("token_id=?", tk.Id).Find(&logs).Error
//	} else {
//		err = LOG_DB.Joins("left join tokens on tokens.id = logs.token_id").Where("tokens.key = ?", strings.TrimPrefix(key, "sk-")).Find(&logs).Error
//	}
//	formatUserLogDetails(logs)
//	return logs, err
//}
//
//func RecordLogDetail(userId int, logType int, content string) {
//	if logType == LogDetailTypeConsume && !common.LogConsumeEnabled {
//		return
//	}
//	username, _ := GetUsernameById(userId, false)
//	log := &LogDetail{
//		UserId:    userId,
//		Username:  username,
//		CreatedAt: common.GetTimestamp(),
//		Type:      logType,
//		Content:   content,
//	}
//	err := LOG_DB.Create(log).Error
//	if err != nil {
//		common.SysError("failed to record log: " + err.Error())
//	}
//}
//
//func RecordErrorLogDetail(c *gin.Context, userId int, channelId int, modelName string, tokenName string, content string, tokenId int, useTimeSeconds int,
//	isStream bool, group string, other map[string]interface{}) {
//	common.LogInfo(c, fmt.Sprintf("record error log: userId=%d, channelId=%d, modelName=%s, tokenName=%s, content=%s", userId, channelId, modelName, tokenName, content))
//	username := c.GetString("username")
//	otherStr := common.MapToJsonStr(other)
//	log := &LogDetail{
//		UserId:           userId,
//		Username:         username,
//		CreatedAt:        common.GetTimestamp(),
//		Type:             LogDetailTypeError,
//		Content:          content,
//		PromptTokens:     0,
//		CompletionTokens: 0,
//		TokenName:        tokenName,
//		ModelName:        modelName,
//		Quota:            0,
//		ChannelId:        channelId,
//		TokenId:          tokenId,
//		UseTime:          useTimeSeconds,
//		IsStream:         isStream,
//		Group:            group,
//		Other:            otherStr,
//	}
//	err := LOG_DB.Create(log).Error
//	if err != nil {
//		common.LogError(c, "failed to record log: "+err.Error())
//	}
//}
//
//func RecordConsumeLogDetail(c *gin.Context, userId int, channelId int, promptTokens int, completionTokens int,
//	modelName string, tokenName string, quota int, content string, tokenId int, userQuota int, useTimeSeconds int,
//	isStream bool, group string, other map[string]interface{}) {
//	common.LogInfo(c, fmt.Sprintf("record consume log: userId=%d, 用户调用前余额=%d, channelId=%d, promptTokens=%d, completionTokens=%d, modelName=%s, tokenName=%s, quota=%d, content=%s", userId, userQuota, channelId, promptTokens, completionTokens, modelName, tokenName, quota, content))
//	if !common.LogConsumeEnabled {
//		return
//	}
//	username := c.GetString("username")
//	otherStr := common.MapToJsonStr(other)
//	log := &LogDetail{
//		UserId:           userId,
//		Username:         username,
//		CreatedAt:        common.GetTimestamp(),
//		Type:             LogDetailTypeConsume,
//		Content:          content,
//		PromptTokens:     promptTokens,
//		CompletionTokens: completionTokens,
//		TokenName:        tokenName,
//		ModelName:        modelName,
//		Quota:            quota,
//		ChannelId:        channelId,
//		TokenId:          tokenId,
//		UseTime:          useTimeSeconds,
//		IsStream:         isStream,
//		Group:            group,
//		Other:            otherStr,
//	}
//
//	// 读取响应内容
//	if customWriter, ok := c.Writer.(*common.CustomResponseWriter); ok {
//		// 使用 Bytes() 避免不必要的字符串复制
//		bodyBytes := customWriter.Body.Bytes()
//		fmt.Printf("ccccc: %s\n", bodyBytes)
//	} else {
//		// 处理转换失败的情况，使用日志库记录
//		common.LogError(c, "c.Writer 不是 *common.CustomResponseWriter 类型")
//	}
//
//	err := LOG_DB.Create(log).Error
//	if err != nil {
//		common.LogError(c, "failed to record log: "+err.Error())
//	}
//	if common.DataExportEnabled {
//		gopool.Go(func() {
//			LogQuotaData(userId, username, modelName, quota, common.GetTimestamp(), promptTokens+completionTokens)
//		})
//	}
//}
//
//func GetAllLogDetails(logType int, startTimestamp int64, endTimestamp int64, modelName string, username string, tokenName string, startIdx int, num int, channel int, group string) (logs []*LogDetail, total int64, err error) {
//	var tx *gorm.DB
//	if logType == LogDetailTypeUnknown {
//		tx = LOG_DB
//	} else {
//		tx = LOG_DB.Where("logs.type = ?", logType)
//	}
//
//	if modelName != "" {
//		tx = tx.Where("logs.model_name like ?", modelName)
//	}
//	if username != "" {
//		tx = tx.Where("logs.username = ?", username)
//	}
//	if tokenName != "" {
//		tx = tx.Where("logs.token_name = ?", tokenName)
//	}
//	if startTimestamp != 0 {
//		tx = tx.Where("logs.created_at >= ?", startTimestamp)
//	}
//	if endTimestamp != 0 {
//		tx = tx.Where("logs.created_at <= ?", endTimestamp)
//	}
//	if channel != 0 {
//		tx = tx.Where("logs.channel_id = ?", channel)
//	}
//	if group != "" {
//		tx = tx.Where("logs."+groupCol+" = ?", group)
//	}
//	err = tx.Model(&LogDetail{}).Count(&total).Error
//	if err != nil {
//		return nil, 0, err
//	}
//	err = tx.Order("logs.id desc").Limit(num).Offset(startIdx).Find(&logs).Error
//	if err != nil {
//		return nil, 0, err
//	}
//
//	channelIds := make([]int, 0)
//	channelMap := make(map[int]string)
//	for _, log := range logs {
//		if log.ChannelId != 0 {
//			channelIds = append(channelIds, log.ChannelId)
//		}
//	}
//	if len(channelIds) > 0 {
//		var channels []struct {
//			Id   int    `gorm:"column:id"`
//			Name string `gorm:"column:name"`
//		}
//		if err = DB.Table("channels").Select("id, name").Where("id IN ?", channelIds).Find(&channels).Error; err != nil {
//			return logs, total, err
//		}
//		for _, channel := range channels {
//			channelMap[channel.Id] = channel.Name
//		}
//		for i := range logs {
//			logs[i].ChannelName = channelMap[logs[i].ChannelId]
//		}
//	}
//
//	return logs, total, err
//}
//
//func GetUserLogDetails(userId int, logType int, startTimestamp int64, endTimestamp int64, modelName string, tokenName string, startIdx int, num int, group string) (logs []*LogDetail, total int64, err error) {
//	var tx *gorm.DB
//	if logType == LogDetailTypeUnknown {
//		tx = LOG_DB.Where("logs.user_id = ?", userId)
//	} else {
//		tx = LOG_DB.Where("logs.user_id = ? and logs.type = ?", userId, logType)
//	}
//
//	if modelName != "" {
//		tx = tx.Where("logs.model_name like ?", modelName)
//	}
//	if tokenName != "" {
//		tx = tx.Where("logs.token_name = ?", tokenName)
//	}
//	if startTimestamp != 0 {
//		tx = tx.Where("logs.created_at >= ?", startTimestamp)
//	}
//	if endTimestamp != 0 {
//		tx = tx.Where("logs.created_at <= ?", endTimestamp)
//	}
//	if group != "" {
//		tx = tx.Where("logs."+groupCol+" = ?", group)
//	}
//	err = tx.Model(&LogDetail{}).Count(&total).Error
//	if err != nil {
//		return nil, 0, err
//	}
//	err = tx.Order("logs.id desc").Limit(num).Offset(startIdx).Find(&logs).Error
//	if err != nil {
//		return nil, 0, err
//	}
//
//	formatUserLogDetails(logs)
//	return logs, total, err
//}
//
//func SearchAllLogDetails(keyword string) (logs []*LogDetail, err error) {
//	err = LOG_DB.Where("type = ? or content LIKE ?", keyword, keyword+"%").Order("id desc").Limit(common.MaxRecentItems).Find(&logs).Error
//	return logs, err
//}
//
//func SearchUserLogDetails(userId int, keyword string) (logs []*LogDetail, err error) {
//	err = LOG_DB.Where("user_id = ? and type = ?", userId, keyword).Order("id desc").Limit(common.MaxRecentItems).Find(&logs).Error
//	formatUserLogDetails(logs)
//	return logs, err
//}
//
//func DeleteOldLogDetail(ctx context.Context, targetTimestamp int64, limit int) (int64, error) {
//	var total int64 = 0
//
//	for {
//		if nil != ctx.Err() {
//			return total, ctx.Err()
//		}
//
//		result := LOG_DB.Where("created_at < ?", targetTimestamp).Limit(limit).Delete(&LogDetail{})
//		if nil != result.Error {
//			return total, result.Error
//		}
//
//		total += result.RowsAffected
//
//		if result.RowsAffected < int64(limit) {
//			break
//		}
//	}
//
//	return total, nil
//}
