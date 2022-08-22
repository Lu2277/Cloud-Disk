package helper

import (
	"Cloud-Disk/core/models"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"path"
	"strconv"
	"strings"
	"time"
)

// GetMd5 生成Md5随机码
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

type UserClaim struct {
	ID       int    `json:"id"`
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("Cloud-Disk key")

// GenerateToken 生成token
func GenerateToken(id int, identity string, name string, expiredTime int) (string, error) {
	UserClaim := &UserClaim{
		ID:       id,
		Identity: identity,
		Name:     name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expiredTime)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	//fmt.Println(signedString)
	return tokenString, nil
}

// AnalyseToken 解析token
func AnalyseToken(tokenString string) (*UserClaim, error) {
	userClaim := new(UserClaim)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("token Invalid:%v", err)
	}
	//fmt.Println(userClaim)
	return userClaim, nil
}

// CreateRandCode  生成随机验证码
func CreateRandCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		randNum := rand.Intn(10)
		code = code + strconv.Itoa(randNum)
	}
	return code
}

// SendCode 发送验证码
func SendCode(Email, code string) error {
	em := email.NewEmail()
	// 设置 sender 发送方的邮箱 ， 此处可以填写自己的邮箱
	em.From = "LU <2591134973@qq.com>"
	// 设置 receiver 接收方的邮箱
	em.To = []string{Email}
	// 设置主题
	em.Subject = "验证码"
	// 简单设置文件发送的内容，暂时设置成纯文本
	em.HTML = []byte("本次验证码为：" + code)
	//设置服务器相关的配置
	err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", "2591134973@qq.com", "lchzddahfbaedjhh", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Send Code Successfully ... ")
	return nil
}

// GetUUID 获取uuid:唯一标识码
func GetUUID() string {
	return uuid.NewV4().String()
}

// CosUpload 上传文件到腾讯云
func CosUpload(r *http.Request) (string, error) {
	//u, _ := url.Parse("https://kk-1313332829.cos.ap-guangzhou.myqcloud.com")
	//b := &cos.BaseURL{BucketURL: u}
	//c := cos.NewClient(b, &http.Client{
	//	Transport: &cos.AuthorizationTransport{
	//		SecretID:  "AKID73j9JDrFMDpBH7h26uTmpcASUGm34pKb",
	//		SecretKey: "3tToukKejX5NQDk13qODm6AyXqthFYcQ",
	//	},
	//})
	file, fileHeader, err := r.FormFile("file")
	name := "cloud-disk/" + GetUUID() + path.Ext(fileHeader.Filename)
	_, err = models.CosClient.Object.Put(context.Background(), name, file, nil)
	if err != nil {
		panic(err)
	}
	//返回可直接访问的路径
	return "https://kk-1313332829.cos.ap-guangzhou.myqcloud.com/" + name, nil
}

// CosFileChuckPrepare 初始化分片上传 返回key和uploadId
func CosFileChuckPrepare(ext string) (string, string, error) {
	key := "cloud-disk/" + GetUUID() + ext
	v, _, err := models.CosClient.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", err
	}
	uploadId := v.UploadID
	return key, uploadId, err
}

// CosFileChuck 分片上传 返回各分片的md5
func CosFileChuck(r *http.Request) (string, error) {
	key := r.PostForm.Get("key")
	UploadID := r.PostForm.Get("upload_id")
	partNumber, err := strconv.Atoi(r.PostForm.Get("part_number"))
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	resp, err := models.CosClient.Object.UploadPart(
		context.Background(), key, UploadID, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", err
	}
	PartETag := strings.Trim(resp.Header.Get("ETag"), "\"") //去掉\和"分隔符
	return PartETag, nil

}

// CosCompleteChuck 分片上传完成
func CosCompleteChuck(key, UploadID string, co []cos.Object) error {
	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, co...)
	_, _, err := models.CosClient.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	return err

}
