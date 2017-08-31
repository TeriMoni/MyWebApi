package controllers

import (
	"MyWebApi/models"
	"MyWebApi/services"
	"MyWebApi/utils"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

type ApiController struct {
	beego.Controller
}

func (this *ApiController) Post() {
	//获取请求的ip地址
	req := this.Ctx.Request
	beego.Info(req.Header.Get("user-agent"))
	addr := req.RemoteAddr
	beego.Info("addr: ", addr)
	//获取请求的数据
	dataList := this.GetString("data")
	fmt.Println(dataList)
	var appdata models.AppData
	if err := json.Unmarshal([]byte(dataList), &appdata); err != nil {
		beego.Error("json string convert error", err)
	}
	var signature string = appdata.Signature //签名字符串
	var userId string = appdata.UserId       //用户id
	var timestamp string = appdata.Timestamp //时间戳
	userKey := utils.GetCache("mongo_user_id")
	if userKey == "" {
		userKey = services.GetCustomerKey(userId)
		utils.SetCache("mongo_user_id", userKey)
	}
	var newSignatureStr string = "timestamp=" + timestamp + "&userId=" + userId + "&userKey=" + userKey
	newSignatureStr = utils.GetMd5(newSignatureStr)
	fmt.Println("timestamp=" + timestamp + "&userId=" + userId + "&userKey=" + userKey)
	fmt.Println("signature:", signature)
	fmt.Println("newSignatureStr:", newSignatureStr)
	if signature != newSignatureStr {
		fmt.Println("签名不合法")
		baseResult := models.BaseResult{1001, "[]", "签名不合法"}
		result, _ := json.Marshal(&baseResult)
		this.Data["json"] = string(result)
		this.ServeJSON()
		return
	}
	baseResult := services.DealAttackInfo([]byte(dataList), addr)
	fmt.Println("the result is :", baseResult)
	this.Data["json"] = baseResult
	this.ServeJSON()
}
