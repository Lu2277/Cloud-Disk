package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var (
	secretID  = "AKID73j9JDrFMDpBH7h26uTmpcASUGm34pKb"
	secretKey = "3tToukKejX5NQDk13qODm6AyXqthFYcQ"
)

//文件上传
func TestCosUpload(t *testing.T) {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://kk-1313332829.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	name := "cloud-disk/123.png"
	// 1.通过字符串上传对象
	//f := strings.NewReader("test")
	//
	//_, err := c.Object.Put(context.Background(), name, f, nil)
	//if err != nil {
	//	panic(err)
	//}
	// 2.通过本地文件上传对象
	_, err := c.Object.PutFromFile(context.Background(), name, "./img/屏幕截图 2022-04-04 220423.png", nil)
	if err != nil {
		panic(err)
	}
	// 3.通过文件流上传对象
	file, err := os.Open("./img/屏幕截图 2022-04-04 220423.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = c.Object.Put(context.Background(), name, file, nil)
	if err != nil {
		panic(err)
	}
}

//分片上传初始化 获取标识分片上传的UploadID
func TestChuckPrepare(t *testing.T) {
	u, _ := url.Parse("https://kk-1313332829.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	name := "cloud-disk/example-123.mp4"
	// 可选opt,如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		t.Fatal(err)
	}
	UploadID := v.UploadID
	fmt.Println(UploadID)
}

//上传分片
func TestChuckUpload(t *testing.T) {
	u, _ := url.Parse("https://kk-1313332829.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	key := "cloud-disk/example-123.mp4"
	UploadID := "1660985979718f7c7ccfaca55ba8f508904ff34b949aef36a72e6c3d906fcc6fcba909fd71"
	f1, err := os.ReadFile("0.chuck")
	f2, err := os.ReadFile("1.chuck")
	resp1, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 1, bytes.NewReader(f1), nil,
	)
	resp2, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 2, bytes.NewReader(f2), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	PartETag1 := resp1.Header.Get("ETag") //ETag即分片文件的Md5值
	PartETag2 := resp2.Header.Get("ETag") //ETag即分片文件的Md5值
	fmt.Println(PartETag1)
	fmt.Println(PartETag2)
}

//完成分片上传
func TestCompletehChuckUpload(t *testing.T) {
	u, _ := url.Parse("https://kk-1313332829.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	key := "cloud-disk/example-123.mp4"
	UploadID := "1660985979718f7c7ccfaca55ba8f508904ff34b949aef36a72e6c3d906fcc6fcba909fd71"
	PartETag1 := "afab97f112e50cac3eaa78c223faaae8"
	PartETag2 := "4ce9e8a09f316186556d73746646351c"
	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cos.Object{PartNumber: 1, ETag: PartETag1},
		cos.Object{PartNumber: 2, ETag: PartETag2},
	)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)

	if err != nil {
		t.Fatal(err)
	}
}
