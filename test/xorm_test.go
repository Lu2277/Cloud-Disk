package test

import (
	"Cloud-Disk/core/models"
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"xorm.io/xorm"
)

//var engine *xorm.Engine

func TestXorm(t *testing.T) {
	//var err error
	engine, err := xorm.NewEngine("mysql", "root:123456@/cloud-disk?charset=utf8mb4")
	if err != nil {
		t.Fatal(err)
	}
	date := make([]*models.User, 0)
	err = engine.Find(&date)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(date)
	if err != nil {
		t.Fatal(err)
	}
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "	")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dst.String())
}
