# Cloud-Disk
基于Go-zero、Xorm开发的网盘系统

## 创建API服务
goctl api new core

## 使用api文件生成代码
goctl api go -api core.api -dir . -style go_zer

## 启动服务
go run core.go -f etc/core-api.yaml
