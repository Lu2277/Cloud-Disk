package test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"
)

//文件的分片
const chuckSize = 10 * 1024 * 1024 //10M 每块分片的大小
func TestGenerateChuckFile(t *testing.T) {
	fileInfo, err := os.Stat("fayewong.mp4")
	if err != nil {
		t.Fatal(err)
	}
	//分片的数量 向上取整数如38/10=3.8，取4
	chuckNum := math.Ceil(float64(fileInfo.Size()) / float64(chuckSize))
	//fmt.Println(chuckNum)
	myFile, err := os.OpenFile("fayewong.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, chuckSize) // 分配空间，存放每一次的分片数据
	for i := 0; i < int(chuckNum); i++ {
		//指定读取文件的起始位置
		myFile.Seek(int64(i*chuckSize), 0)
		//读取到最后一块分片  如150M分成100M、50M两个分片，最后一块分片<=分片的大小
		if chuckSize > fileInfo.Size()-int64(i*chuckSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chuckSize))
		}
		myFile.Read(b)
		file, err := os.OpenFile("./"+strconv.Itoa(i)+".chuck", os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			t.Fatal(err)
		}
		file.Write(b)
		file.Close()
	}
	myFile.Close()
}

//分片文件的合并
func TestMergeChuckFile(t *testing.T) {
	myFile, err := os.OpenFile("fayewong-2.mp4", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat("fayewong.mp4")
	if err != nil {
		t.Fatal(err)
	}
	// 计算分片个数, 正常应该由前端传来, 这里测试时自行计算
	chuckNum := math.Ceil(float64(fileInfo.Size()) / float64(chuckSize))
	for i := 0; i < int(chuckNum); i++ {
		file, err := os.OpenFile("./"+strconv.Itoa(i)+".chuck", os.O_RDONLY, 0777)
		if err != nil {
			t.Fatal(err)
		}
		all, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		myFile.Write(all)
		file.Close()
	}
	myFile.Close()
}

//文件的一致性校验
func TestCheckMd5(t *testing.T) {
	//获取文件的信息
	file1, err := os.OpenFile("fayewong.mp4", os.O_RDONLY, 0666)
	file2, err := os.OpenFile("fayewong-2.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b1, err := ioutil.ReadAll(file1)
	b2, err := ioutil.ReadAll(file2)
	if err != nil {
		t.Fatal()
	}
	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s1 == s2)

}
