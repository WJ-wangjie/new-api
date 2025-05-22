package model

import (
	"fmt"
	"gorm.io/gorm"
	"one-api/common"
	"sync"
	"time"
)

// QuotaData 柱状图数据
type QuotaData struct {
	Id          int    `json:"id"`
	UserID      int    `json:"user_id" gorm:"index"`
	Username    string `json:"username" gorm:"index:idx_qdt_model_user_name,priority:2;size:64;default:''"`
	ChannelId   int    `json:"channel_id"`
	ModelName   string `json:"model_name" gorm:"index:idx_qdt_model_user_name,priority:1;size:64;default:''"`
	CreatedAt   int64  `json:"created_at" gorm:"bigint;index:idx_qdt_created_at,priority:2"`
	TokenUsed   int    `json:"token_used" gorm:"default:0"`
	Count       int    `json:"count" gorm:"default:0"`
	Quota       int    `json:"quota" gorm:"default:0"`
	ChannelName string `json:"channel_name" gorm:"-"`
}

// QuotaDownloadData 下载配额数据结构体
type QuotaDownloadData struct {
	Date        string  `json:"date"`
	ChannelName string  `json:"channel_name"`
	ModelName   string  `json:"model_name"`
	Count       int     `json:"count"`
	Quota       int     `json:"quota"`
	Amount      float64 `json:"amount"`
	TokenUsed   int     `json:"token_used"`
}

func UpdateQuotaData() {
	// recover
	defer func() {
		if r := recover(); r != nil {
			common.SysLog(fmt.Sprintf("UpdateQuotaData panic: %s", r))
		}
	}()
	for {
		if common.DataExportEnabled {
			common.SysLog("正在更新数据看板数据...")
			SaveQuotaDataCache()
		}
		time.Sleep(time.Duration(common.DataExportInterval) * time.Minute)
	}
}

var CacheQuotaData = make(map[string]*QuotaData)
var CacheQuotaDataLock = sync.Mutex{}

func logQuotaDataCache(userId int, username string, channelId int, modelName string, quota int, createdAt int64, tokenUsed int) {
	key := fmt.Sprintf("%d-%s--%d-%s-%d", userId, username, channelId, modelName, createdAt)
	quotaData, ok := CacheQuotaData[key]
	if ok {
		quotaData.Count += 1
		quotaData.Quota += quota
		quotaData.TokenUsed += tokenUsed
	} else {
		quotaData = &QuotaData{
			UserID:    userId,
			Username:  username,
			ModelName: modelName,
			ChannelId: channelId,
			CreatedAt: createdAt,
			Count:     1,
			Quota:     quota,
			TokenUsed: tokenUsed,
		}
	}
	CacheQuotaData[key] = quotaData
}

func LogQuotaData(userId int, username string, channelId int, modelName string, quota int, createdAt int64, tokenUsed int) {
	// 只精确到小时
	createdAt = createdAt - (createdAt % 3600)

	CacheQuotaDataLock.Lock()
	defer CacheQuotaDataLock.Unlock()
	logQuotaDataCache(userId, username, channelId, modelName, quota, createdAt, tokenUsed)
}

func SaveQuotaDataCache() {
	CacheQuotaDataLock.Lock()
	defer CacheQuotaDataLock.Unlock()
	size := len(CacheQuotaData)
	// 如果缓存中有数据，就保存到数据库中
	// 1. 先查询数据库中是否有数据
	// 2. 如果有数据，就更新数据
	// 3. 如果没有数据，就插入数据
	for _, quotaData := range CacheQuotaData {
		quotaDataDB := &QuotaData{}
		DB.Table("quota_data").Where("user_id = ? and username = ? and channel_id = ? and model_name = ? and created_at = ?",
			quotaData.UserID, quotaData.Username, quotaData.ChannelId, quotaData.ModelName, quotaData.CreatedAt).First(quotaDataDB)
		if quotaDataDB.Id > 0 {
			//quotaDataDB.Count += quotaData.Count
			//quotaDataDB.Quota += quotaData.Quota
			//DB.Table("quota_data").Save(quotaDataDB)
			increaseQuotaData(quotaData.UserID, quotaData.Username, quotaData.ChannelId, quotaData.ModelName, quotaData.Count, quotaData.Quota, quotaData.CreatedAt, quotaData.TokenUsed)
		} else {
			DB.Table("quota_data").Create(quotaData)
		}
	}
	CacheQuotaData = make(map[string]*QuotaData)
	common.SysLog(fmt.Sprintf("保存数据看板数据成功，共保存%d条数据", size))
}

func increaseQuotaData(userId int, username string, channelId int, modelName string, count int, quota int, createdAt int64, tokenUsed int) {
	err := DB.Table("quota_data").Where("user_id = ? and username = ? and channel_id = ? and model_name = ? and created_at = ?",
		userId, username, channelId, modelName, createdAt).Updates(map[string]interface{}{
		"count":      gorm.Expr("count + ?", count),
		"quota":      gorm.Expr("quota + ?", quota),
		"token_used": gorm.Expr("token_used + ?", tokenUsed),
	}).Error
	if err != nil {
		common.SysLog(fmt.Sprintf("increaseQuotaData error: %s", err))
	}
}

// 提取公共查询函数
func queryQuotaData(db *gorm.DB) (quotaData []*QuotaData, err error) {
	// 定义匿名结构体用于接收查询结果
	var results []struct {
		QuotaData
		ChannelName string `json:"channel_name"`
	}
	err = db.Find(&results).Error

	if err != nil {
		return nil, err
	}
	// 将结果转换为 QuotaData 切片
	quotaData = make([]*QuotaData, len(results))
	for i, item := range results {
		quotaData[i] = &item.QuotaData
		quotaData[i].ChannelName = item.ChannelName
	}
	return quotaData, err
}

func GetQuotaDataByUsername(username string, startTime int64, endTime int64) (quotaData []*QuotaData, err error) {
	// 从quota_data表中查询数据
	db := DB.Table("quota_data").
		Select("channels.name as channel_name,model_name, sum(count) as count, sum(quota) as quota, sum(token_used) as token_used, created_at").
		Joins("LEFT JOIN channels ON quota_data.channel_id = channels.id").
		Where("username = ? and created_at >= ? and created_at <= ?", username, startTime, endTime).
		Group("model_name, created_at, channels.name")
	return queryQuotaData(db)
}

func GetQuotaDataByUserId(userId int, startTime int64, endTime int64) (quotaData []*QuotaData, err error) {
	// 从quota_data表中查询数据
	//err = DB.Table("quota_data").Where("user_id = ? and created_at >= ? and created_at <= ?", userId, startTime, endTime).Find(&quotaDatas).Error

	db := DB.Table("quota_data").
		Select("channels.name as channel_name,model_name, sum(count) as count, sum(quota) as quota, sum(token_used) as token_used, created_at").
		Joins("LEFT JOIN channels ON quota_data.channel_id = channels.id").
		Where("user_id = ? and created_at >= ? and created_at <= ?", userId, startTime, endTime).
		Group("model_name, created_at, channels.name")

	return queryQuotaData(db)
}

func GetAllQuotaDates(startTime int64, endTime int64, username string) (quotaData []*QuotaData, err error) {
	if username != "" {
		return GetQuotaDataByUsername(username, startTime, endTime)
	}
	// 从quota_data表中查询数据
	// only select model_name, sum(count) as count, sum(quota) as quota, model_name, created_at from quota_data group by model_name, created_at;
	//err = DB.Table("quota_data").Where("created_at >= ? and created_at <= ?", startTime, endTime).Find(&quotaDatas).Error
	db := DB.Table("quota_data").Select("channels.name as channel_name,model_name, sum(count) as count, sum(quota) as quota, sum(token_used) as token_used, created_at").
		Joins("LEFT JOIN channels ON quota_data.channel_id = channels.id").
		Where("created_at >= ? and created_at <= ?", startTime, endTime).
		Group("model_name, created_at, channels.name")

	return queryQuotaData(db)
}

func DownloadQuotaData(userId int, username string, startTime int64, endTime int64) (quotaDownloadData []*QuotaDownloadData, err error) {
	// 构建递归 CTE 日期序列
	dateSeriesCTE := `
	WITH RECURSIVE filtered_data AS (
		SELECT *
		FROM quota_data
		WHERE 1=1 `

	var args []interface{}
	if userId != 0 {
		dateSeriesCTE += " AND user_id = ?"
		args = append(args, userId)
	}
	if username != "" {
		dateSeriesCTE += " AND username = ?"
		args = append(args, username)
	}
	if startTime != 0 {
		dateSeriesCTE += " AND created_at >= ?"
		args = append(args, startTime)
	}
	if endTime != 0 {
		dateSeriesCTE += " AND created_at <= ?"
		args = append(args, endTime)
	}

	dateSeriesCTE = dateSeriesCTE + ` 
	),
	date_range AS (
	SELECT
			MIN(DATE(FROM_UNIXTIME(created_at))) AS min_date,
			MAX(DATE(FROM_UNIXTIME(created_at))) AS max_date
		FROM filtered_data
	),
	date_series AS (
		SELECT min_date AS date FROM date_range
		UNION ALL
		SELECT DATE_ADD(date, INTERVAL 1 DAY)
		FROM date_series, date_range
		WHERE date < max_date
	),
	channel_models AS (
		SELECT DISTINCT
			c.name AS channel_name,
			q.model_name
		FROM filtered_data q
		JOIN channels c ON q.channel_id = c.id
		WHERE q.model_name IS NOT NULL
	),
	daily_stats AS (
		SELECT
			DATE(FROM_UNIXTIME(q.created_at)) AS date,
			c.name AS channel_name,
			q.model_name,
			SUM(q.count) AS count,
			SUM(q.quota) AS quota,
			SUM(q.token_used) AS token_used
		FROM filtered_data q
		JOIN channels c ON q.channel_id = c.id
		WHERE q.model_name IS NOT NULL
		GROUP BY DATE(FROM_UNIXTIME(q.created_at)), c.name, q.model_name
	)
	SELECT
		cm.channel_name,
		cm.model_name,
		DATE_FORMAT(date_series.date, '%Y-%m-%d') AS date, -- 使用 DATE_FORMAT 格式化日期
		COALESCE(ds.count, 0) AS count,
		COALESCE(ds.quota, 0) AS quota,
		COALESCE(ds.quota, 0) / 500000.0 AS amount,
		COALESCE(ds.token_used, 0) AS token_used
	FROM date_series
	CROSS JOIN channel_models cm
	LEFT JOIN daily_stats ds ON
		date_series.date = ds.date AND
		cm.channel_name = ds.channel_name AND
		cm.model_name = ds.model_name
	ORDER BY cm.channel_name, cm.model_name, date_series.date;
	`

	// 执行查询
	err = DB.Raw(dateSeriesCTE, args...).Scan(&quotaDownloadData).Error
	if err != nil {
		return nil, err
	}
	return quotaDownloadData, nil
}
