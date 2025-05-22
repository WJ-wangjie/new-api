package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
	"net/http"
	"one-api/model"
	"strconv"
	"time"
)

func GetAllQuotaDates(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	dates, err := model.GetAllQuotaDates(startTimestamp, endTimestamp, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    dates,
	})
	return
}

func GetUserQuotaDates(c *gin.Context) {
	userId := c.GetInt("id")
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	// 判断时间跨度是否超过 1 个月
	if endTimestamp-startTimestamp > 2592000 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "时间跨度不能超过 1 个月",
		})
		return
	}
	dates, err := model.GetQuotaDataByUserId(userId, startTimestamp, endTimestamp)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    dates,
	})
	return
}

func DownloadQuotaDates(c *gin.Context) {
	userId := c.GetInt("id")
	username := c.Query("username")
	startTimestamp, err := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "start_timestamp 参数解析失败",
		})
		return
	}
	endTimestamp, err := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "end_timestamp 参数解析失败",
		})
		return
	}

	// 调用模型层函数获取下载数据
	dates, err := model.DownloadQuotaData(userId, username, startTimestamp, endTimestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 创建一个新的 Excel 文件
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("对账单数据")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建 Excel 工作表失败",
		})
		return
	}

	// 添加表头，这里假设 dates 是一个结构体切片，你需要根据实际情况调整
	headerRow := sheet.AddRow()
	// 示例表头，根据实际数据结构修改
	headerCells := []string{"日期", "渠道", "模型", "次数", "quota", "金额", "token使用"}
	for _, header := range headerCells {
		cell := headerRow.AddCell()
		cell.SetValue(header)
	}

	// 填充数据，根据实际数据结构修改
	for _, data := range dates {
		row := sheet.AddRow()
		// 示例数据填充，根据实际数据结构修改
		cell := row.AddCell()
		cell.SetValue(data.Date)
		cell = row.AddCell()
		cell.SetValue(data.ChannelName)
		cell = row.AddCell()
		cell.SetValue(data.ModelName)
		cell = row.AddCell()
		cell.SetValue(data.Count)
		cell = row.AddCell()
		cell.SetValue(data.Quota)
		cell = row.AddCell()
		cell.SetValue(data.Amount)
		cell = row.AddCell()
		cell.SetValue(data.TokenUsed)
	}

	// 创建一个 bytes.Buffer 用于存储 Excel 数据
	var buf bytes.Buffer
	// 使用 Write 方法将 Excel 文件内容写入 bytes.Buffer
	err = file.Write(&buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "写入 Excel 数据失败",
		})
		return
	}
	// 将 bytes.Buffer 中的数据转换为字节切片
	excelData := buf.Bytes()

	// 设置响应头，用于文件下载
	c.Header("Content-Disposition", "attachment; filename=对账单-"+time.Now().Format("20060102150405")+".xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(excelData)))

	// 写入文件内容
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
	return
}
