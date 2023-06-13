# SAP 环境自动配置

本程序只支持在 windows 环境下 SAP 环境自动设置

## 安装

请安装 golang 语言执行程序 windows 版本。

https://go.dev/dl/

下载本项目源代码。

```
git clone https://github.com/wwsheng009/sap_check_config.git
cd sap_check_config
go mod tidy
```

## 编译

```shell

go get github.com/akavel/rsrc
# go 1.8
go install github.com/akavel/rsrc@latest

rsrc -manifest .\main.exe.manifest -ico ./assets/main.ico -o SAP环境自动配置.syso

rsrc -manifest .\main_admin.exe.manifest -ico ./assets/main.ico -o SAP环境自动配置.syso
go build -ldflags="-H windowsgui" -o SAP环境自动配置.exe
```
