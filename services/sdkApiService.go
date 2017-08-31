package services

import (
	"MyWebApi/models"
	"MyWebApi/utils"
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoUrl = "212.188.172.216:27017"
	unExist  = "notExist"
	redisKey = "test_key"
)

var (
	mgoSession *mgo.Session
	dataBase   = "threatPerception"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mongoUrl)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

func GetCustomerKey(userId string) string {
	result := models.Customer{}
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(&result)
	}
	err := witchCollection("customer_key", query)
	if err != nil {
		//		panic(err)
		return unExist
	}
	fmt.Println("Customer:", result.UserId, result.UserKey, result.CreateTime)
	return result.UserKey
}

func DealAttackInfo(param []byte, ip string) models.BaseResult {

	mongList := list.New() //定义一个集合用来存放攻击结果集，存入redis
	var appData models.AppData
	json.Unmarshal(param, &appData)
	//	data, _ := json.Marshal(appData)
	currentAppInfo := appData.CurrentAppInfo
	networkInfo := appData.NetworkInfo
	systemInfo := appData.SystemInfo
	attackMonitor := appData.AttackMonitor
	var appName string = currentAppInfo.AppName
	fmt.Println("攻击检测app名称:", appName)
	//处理具体业务逻辑
	var customerApp = models.App{}
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"appName": appName}).One(&customerApp)
	}
	err := witchCollection("customer_app", query)
	if err != nil {
		panic(err)
	}
	fmt.Println("攻击检测的客户id:", customerApp.UserId)
	var attackResult = models.AttackResult{}
	//判断是否为debug 调试
	var isDebug bool = attackMonitor.IsDebug
	var dexHook string = attackMonitor.DexHook
	var soHook string = attackMonitor.SoHook
	var rePackage string = attackMonitor.RePackage
	var hijackInfo models.HijackInfo = attackMonitor.IsHijack
	if isDebug {
		fmt.Println("攻击类型为debug调试！")
		attackResult.IsDebug = true
		insertValue := getInsertValue(customerApp, currentAppInfo, networkInfo, systemInfo, appData, ip)
		mongList.PushBack(insertValue)
	}
	if dexHook != "" {
		for _, info := range utils.GetAttackNames() {
			if info == dexHook {
				fmt.Println("攻击类型为dexHook攻击！,攻击类型:" + info)
				attackResult.DexHook = true
				insertValue := getInsertValue(customerApp, currentAppInfo, networkInfo, systemInfo, appData, ip)
				insertValue["attackType"] = "dexHook"
				mongList.PushBack(insertValue)
				break
			}
		}
	}
	if soHook != "" {
		for _, info := range utils.GetAttackNames() {
			if info == soHook {
				fmt.Println("攻击类型为soHook攻击！,攻击类型:" + info)
				attackResult.SoHook = true
				insertValue := getInsertValue(customerApp, currentAppInfo, networkInfo, systemInfo, appData, ip)
				insertValue["attackType"] = "soHook"
				mongList.PushBack(insertValue)
				break
			}
		}
	}

	if rePackage != "" {
		flag := IsRePackageApp(currentAppInfo.MD5)
		if flag {
			fmt.Println("攻击类型为二次打包！")
			attackResult.RePackage = true
			insertValue := getInsertValue(customerApp, currentAppInfo, networkInfo, systemInfo, appData, ip)
			insertValue["attackType"] = "rePackage"
			mongList.PushBack(insertValue)
		}

	}
	hijackpackage := hijackInfo.AppPackage
	if hijackpackage != "" {
		for _, info := range utils.GetDangerPackages() {
			if info == hijackpackage {
				fmt.Println("攻击类型为isHijack攻击")
				attackResult.IsHijack = true
				insertValue := getInsertValue(customerApp, currentAppInfo, networkInfo, systemInfo, appData, ip)
				insertValue["attackType"] = "isHijack"
				mongList.PushBack(insertValue)
				break
			}
		}
	}
	//存mongoList到redis队列，返回结果集
	for p := mongList.Front(); p != nil; p = p.Next() {
		data, _ := json.Marshal(p.Value)
		utils.Push(redisKey, string(data))
	}
	baseResult := models.BaseResult{0, attackResult, "检测完成"}
	//转json字符串输出
	//	result, _ := json.Marshal(&baseResult)
	return baseResult
}

//查询库中是否存在，判断是否是二次打包
func IsRePackageApp(appMd5 string) bool {
	var customerApp = models.App{}
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"appMd5": appMd5}).One(&customerApp)
	}
	err := witchCollection("customer_app", query)
	if err != nil {
		//not found
		return false
	}
	return true
}

//format结果集
func getInsertValue(customerApp models.App, currentAppInfo models.CurrentAppInfo, networkInfo models.NetWorkInfo, systemInfo models.SystemInfo, appData models.AppData, ip string) map[string]interface{} {
	resultMap := make(map[string]interface{})
	resultMap["userId"] = customerApp.UserId
	resultMap["attackTime"] = time.Now().Format("2006-01-02 15:04:05")
	resultMap["phoneNumber"] = networkInfo.PhoneNumber
	resultMap["deviceId"] = getDeviceId(networkInfo)
	resultMap["ip"] = ip
	resultMap["operatingSystem"] = appData.Type
	resultMap["mobileBrand"] = systemInfo.Brand
	resultMap["networkType"] = networkInfo.NetWorkType
	resultMap["mac"] = networkInfo.MAC
	resultMap["IMEI"] = networkInfo.IMEI
	resultMap["IMSI"] = networkInfo.IMSI
	resultMap["timeZone"] = systemInfo.TimeZone
	resultMap["location"] = systemInfo.Location
	//通过请求接口获取地理位置信息
	addressMap := getAddressFromIp(ip, "GBK")
	if addressMap["province"] != nil {
		resultMap["attribution"] = addressMap["province"]
	} else {
		resultMap["attribution"] = ""
	}
	if addressMap["city"] != nil {
		resultMap["city"] = addressMap["city"]
	} else {
		resultMap["city"] = ""
	}
	resultMap["appMd5"] = currentAppInfo.MD5
	resultMap["appVersion"] = currentAppInfo.VersionName
	resultMap["appName"] = currentAppInfo.AppName
	resultMap["packageName"] = currentAppInfo.PackageName
	installTime, err := strconv.Atoi(currentAppInfo.FirstInstallTime)
	if err != nil {
		log.Fatal(err)
	}
	updateTime, err := strconv.Atoi(currentAppInfo.LastUpdateTime)
	if err != nil {
		log.Fatal(err)
	}
	//格式化时间毫秒数为时间
	resultMap["installTime"] = time.Unix(int64(installTime)/1000, 0).Format("2006-01-02 15:04:05")
	resultMap["updateTime"] = time.Unix(int64(updateTime)/1000, 0).Format("2006-01-02 15:04:05")
	resultMap["digitalCertificate"] = currentAppInfo.SerialNumber
	resultMap["systemVersion"] = systemInfo.Android_version
	return resultMap
}

//获取设备id做md5处理
func getDeviceId(networkInfo models.NetWorkInfo) string {
	return utils.GetMd5(networkInfo.IMEI + networkInfo.IMSI + networkInfo.MAC)
}

//通过请求网易接口获取ip地址对应的地域信息
func getAddressFromIp(ip string, enCoding string) map[string]interface{} {
	var result = make(map[string]interface{}, 0)
	postUrl := "http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=json&ip=" + "221.228.46.31"
	resp, err := http.Get(postUrl)
	if err != nil {
		panic(err)
		fmt.Println("请求获取地域信息接口错误")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
	}
	return result
}

//unicode转中文
func tansferCode(source string) string {

	sUnicodev := strings.Split(source, "\\u")
	var context string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)
	}
	fmt.Println(context)
	return context
}
