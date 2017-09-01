package utils

//redis 工具类
import (
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

var (
	// 定义常量
	RedisClient    *redis.Pool
	REDIS_HOST     string
	REDIS_DB       string
	REDIS_PASSWORD string
)

func init() {
	// 初始化redis客户端地址和数据库读取本地配置文件
	REDIS_HOST = beego.AppConfig.String("redis.host")
	REDIS_DB = beego.AppConfig.String("redis.db")
	REDIS_PASSWORD = beego.AppConfig.String("redis.password")
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     beego.AppConfig.DefaultInt("redis.maxidle", 3),
		MaxActive:   beego.AppConfig.DefaultInt("redis.maxactive", 30),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				beego.Error("Connect to redis error", err)
				return nil, err
			}
			if _, err := c.Do("AUTH", REDIS_PASSWORD); err != nil {
				c.Close()
				beego.Error("password is error", err)
				return nil, err
			}
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}

//取缓存
func GetCache(key string) string {
	// 从池里获取连接
	rs := RedisClient.Get()
	// 用完后将连接放回连接池
	defer rs.Close()
	value, err := redis.String(rs.Do("GET", key))
	if err != nil {
		beego.Error("redis get cache error,key:", key)
		return ""
	}
	return value
}

//设置缓存
func SetCache(key string, value string) {

	rs := RedisClient.Get()
	defer rs.Close()
	_, err := rs.Do("SET", key, value)
	if err != nil {
		beego.Error("redis set cache failed:", err)
	}
}

//设置带过期时间的缓存
func SetCacheAndTime(key string, value string, expireTime int) {

	rs := RedisClient.Get()
	defer rs.Close()
	_, err := rs.Do("SET", key, value)
	if expireTime > 0 {
		rs.Do("EXPIRE", key, expireTime)
	}
	if err != nil {
		beego.Error("redis set cache and expire time failed:", err)
	}
}

//删除缓存
func DeleteCache(key string) bool {
	rs := RedisClient.Get()
	defer rs.Close()
	_, err := rs.Do("DEL", key)
	if err != nil {
		beego.Error("redis delelte cache failed:", err)
		return false
	}
	return true
}

//入队列
func Push(key string, value string) {
	rs := RedisClient.Get()
	defer rs.Close()
	_, err := rs.Do("rpush", key, value)
	if err != nil {
		beego.Error("redis push cache queue failed:", err)
	}
}

//出队列
func Pop(key string) {
	rs := RedisClient.Get()
	defer rs.Close()
	_, err := rs.Do("lpop", key)
	if err != nil {
		beego.Error("redis pop cache queue failed:", err)
	}
}

//获取队列长度
func Length(key string) int64 {
	rs := RedisClient.Get()
	defer rs.Close()
	lenth, err := rs.Do("llen", key)
	if err != nil {
		beego.Error("redis get quenue length failed:", err)
	}
	count, ok := lenth.(int64)
	if !ok {
		beego.Error("类型转换错误!", err)
		return 0
	}
	return count
}

func IsExist(key string) bool {
	rs := RedisClient.Get()
	defer rs.Close()
	is_key_exit, err := redis.Bool(rs.Do("EXISTS", key))
	if err != nil {
		beego.Error("redis get key failed:", err)
		return false
	} else {
		return is_key_exit
	}
}
