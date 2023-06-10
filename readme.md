## icon

```shell

go get github.com/akavel/rsrc
# go 1.8
go install github.com/akavel/rsrc@latest

rsrc -manifest .\main.exe.manifest -ico ./main.ico -o sap_check.syso
go build -ldflags="-H windowsgui" -o sap_check.exe
```
