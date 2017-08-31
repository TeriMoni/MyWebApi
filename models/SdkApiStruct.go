package models

import (
	"gopkg.in/mgo.v2/bson"
)

type BaseResult struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type AttackResult struct {
	IsDebug   bool `json:"isDebug"`
	SoHook    bool `json:"soHook"`
	DexHook   bool `json:"dexHook"`
	RePackage bool `json:"rePackage"`
	IsHijack  bool `json:"isHijack"`
}

//所有攻击类型
type MonitorInfo struct {
	DexHook   string
	IsDebug   bool
	SoHook    string
	RePackage string
	IsHijack  HijackInfo
}

type HijackInfo struct {
	AppPackage string
	CerFin     string
}

type AppData struct {
	AttackMonitor     MonitorInfo
	Signature         string
	Timestamp         string
	Type              string
	UserId            string
	CurrentAppInfo    CurrentAppInfo
	EnvironmentalInfo EnvironmentInfo
	NetworkInfo       NetWorkInfo
	SystemInfo        SystemInfo
}

type CurrentAppInfo struct {
	MD5              string
	AppName          string
	Brand            string
	FirstInstallTime string
	Issuer           string
	LastModified     string
	LastUpdateTime   string
	PackageName      string
	PubKey           string
	SerialNumber     string
	SignName         string
	SubjectDN        string
	VersionCode      string
	VersionName      string
}

type EnvironmentInfo struct {
	IsEmulator bool
	IsRoot     bool
}

type NetWorkInfo struct {
	IMEI            string
	IMSI            string
	MAC             string
	NetWorkOperator string
	NetWorkType     string
	PhoneNumber     string
	SimNo           string
}

type SystemInfo struct {
	Android_version  string
	Brand            string
	CpuCoreNum       string
	CpuName          string
	Location         string
	Manufacture      string
	MaxCpuFreq       string
	Model            string
	RamMemory        string
	Romsize          int32
	ScreenResolution string
	StartTime        string
	TimeZone         string
}
type Customer struct {
	Id_        bson.ObjectId `bson:"_id"`
	UserId     int32         `bson:"userId"`
	UserKey    string        `bson:"userKey"`
	CreateTime string        `bson:"createTime"`
}

type App struct {
	Id_     bson.ObjectId `bson:"_id"`
	UserId  int32         `bson:"userId"`
	AppName string        `bson:"appName"`
	AppMd5  string        `bson:"appMd5"`
}
