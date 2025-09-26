package service

import (
	"gossh/app/utils"
	"gossh/gin"
	"log/slog"
	"net/http"
)

// UbusAction 用于调用 ubus 并返回结果
func UbusAction(c *gin.Context) {
	// 请求参数结构体
	type ubusRequest struct {
		Service string                 `form:"service" json:"service" binding:"required"`
		Method  string                 `form:"method"  json:"method"  binding:"required"`
		Params  map[string]interface{} `form:"params"  json:"params"`
	}

	var req ubusRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("绑定数据错误", "err_msg", err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "输入数据不合法"})
		return
	}

	// 调用 utils.getDataFromUbus
	result, err := utils.GetDataFromUbus(req.Service, req.Method, req.Params)
	if err != nil {
		slog.Error("Ubus 调用失败", "err_msg", err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "ubus 调用失败", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"msg":    "success",
		"result": result,
	})
}
