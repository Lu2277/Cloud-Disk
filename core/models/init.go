package models

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"net/url"
	"xorm.io/xorm"
)

var Engine = Init()
var RDB = InitRedis()
var CosClient = InitCosClient()

func Init() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:123456@/cloud-disk?charset=utf8mb4")
	if err != nil {
		log.Printf("Xorm New Engine Error:%v", err)
		return nil
	}
	return engine
}
func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}
func InitCosClient() *cos.Client {
	u, _ := url.Parse("https://kk-1313332829.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "AKID73j9JDrFMDpBH7h26uTmpcASUGm34pKb",
			SecretKey: "3tToukKejX5NQDk13qODm6AyXqthFYcQ",
		},
	})
	return c
}
